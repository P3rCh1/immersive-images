import json
import uuid
import os
import sys
from datetime import datetime, timezone
from http.server import BaseHTTPRequestHandler, HTTPServer
from urllib.parse import urlparse, parse_qs

IMAGES = {}         
IMAGE_BYTES = {}    

SCRIPT_DIR = os.path.dirname(os.path.abspath(__file__))
TEST_IMAGE_PATH = os.path.join(SCRIPT_DIR, "test_image.png")

MINI_PNG = bytes.fromhex(
    "89504e470d0a1a0a0000000d4948445200000001000000010806"
    "0000001f15c4890000000d49444154789c6360000000020001e2"
    "21bc330000000049454e44ae426082"
)

def load_test_image():
    if os.path.exists(TEST_IMAGE_PATH):
        with open(TEST_IMAGE_PATH, "rb") as f:
            data = f.read()
        print(f"[mock] Loaded test image: {len(data)} bytes")
        return data
    print("[mock] test_image.png not found, using 1x1 PNG")
    return MINI_PNG

TEST_IMAGE = load_test_image()

def now_iso():
    return datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")

def error_response(msg):
    return {"error": {"message": msg}}

def default_settings():
    return {
        "style": "PATCHES",
        "palette": "NEON",
        "width": 512,
        "height": 512,
        "scale": 1.0,
    }

def make_image_meta(seed, name, settings=None):
    return {
        "uuid": str(uuid.uuid4()),
        "name": name,
        "created_at": now_iso(),
        "seed": seed,
        "settings": settings or default_settings(),
    }

class Handler(BaseHTTPRequestHandler):
    
    def _send_json(self, data, status=200):
        body = json.dumps(data, ensure_ascii=False).encode("utf-8")
        self.send_response(status)
        self.send_header("Content-Type", "application/json; charset=utf-8")
        self.send_header("Content-Length", str(len(body)))
        self.end_headers()
        self.wfile.write(body)
    
    def _send_png(self, png_bytes, status=200):
        self.send_response(status)
        self.send_header("Content-Type", "image/png")
        self.send_header("Content-Length", str(len(png_bytes)))
        self.send_header("Cache-Control", "public, max-age=31536000, immutable")
        self.send_header("Content-Disposition", 'inline; filename="image.png"')
        self.end_headers()
        self.wfile.write(png_bytes)
    
    def _send_no_content(self):
        self.send_response(204)
        self.end_headers()
    
    def _read_json_body(self):
        length = int(self.headers.get("Content-Length", 0))
        if length == 0:
            return None
        raw = self.rfile.read(length).decode("utf-8")
        return json.loads(raw)
    
    def _extract_uuid(self, path, suffix=""):
        prefix = "/images/"
        if not path.startswith(prefix):
            return None
        rest = path[len(prefix):]
        if suffix and rest.endswith(suffix):
            rest = rest[:-len(suffix)]
        return rest.split("/")[0]
    
    def do_GET(self):
        parsed = urlparse(self.path)
        path = parsed.path
        query = parse_qs(parsed.query)
        
        if path == "/healthz":
            return self._send_json({"status": "ok"})
        
        if path == "/images":
            limit = int(query.get("limit", [50])[0])
            start_uuid = query.get("start", [None])[0]
            
            items = sorted(IMAGES.values(), 
                          key=lambda x: x["created_at"], 
                          reverse=True)
            
            if start_uuid:
                idx = next((i for i, x in enumerate(items) 
                           if x["uuid"] == start_uuid), 0)
                items = items[idx:]
            
            page = items[:limit]
            next_uuid = items[limit]["uuid"] if len(items) > limit else None
            
            response = {
                "items": page,
                "total": len(IMAGES),
                "limit": limit,
            }
            if next_uuid:
                response["next"] = next_uuid
            
            return self._send_json(response)
        
        if path.startswith("/images/") and path.endswith("/content"):
            img_uuid = self._extract_uuid(path, "/content")
            if img_uuid not in IMAGES:
                return self._send_json(error_response("Image not found"), 404)
            return self._send_png(IMAGE_BYTES.get(img_uuid, TEST_IMAGE))
        
        if path.startswith("/images/"):
            img_uuid = self._extract_uuid(path)
            if img_uuid not in IMAGES:
                return self._send_json(error_response("Image not found"), 404)
            return self._send_json(IMAGES[img_uuid])
        
        return self._send_json(error_response("Not found"), 404)
    
    def do_POST(self):
        if self.path != "/images":
            return self._send_json(error_response("Not found"), 404)
        
        try:
            body = self._read_json_body()
        except Exception as e:
            return self._send_json(error_response(f"Invalid JSON: {e}"), 400)
        
        if body is None:
            return self._send_json(error_response("Empty body"), 400)
        
        seed = body.get("seed")
        name = (body.get("name") or "").strip()
        settings = body.get("settings")
        
        if seed is None:
            return self._send_json(error_response("seed is required"), 400)
        if not isinstance(seed, int):
            return self._send_json(error_response("seed must be integer"), 400)
        if not name:
            return self._send_json(error_response("name is required"), 400)
        if len(name) > 255:
            return self._send_json(error_response("name too long"), 400)
        
        if settings:
            style = settings.get("style")
            palette = settings.get("palette")
            valid_styles = ["PATCHES", "SPOTS", "RINGS", 
                           "TERRACES", "WAVES", "MOTTLED"]
            valid_palettes = ["NEON", "DARK", "TOXIC", "RAINBOW", 
                             "OCEAN", "LAVA", "MONO"]
            
            if style and style not in valid_styles:
                return self._send_json(
                    error_response(f"Invalid style: {style}"), 400)
            if palette and palette not in valid_palettes:
                return self._send_json(
                    error_response(f"Invalid palette: {palette}"), 400)
        
        meta = make_image_meta(seed, name, settings)
        IMAGES[meta["uuid"]] = meta
        IMAGE_BYTES[meta["uuid"]] = TEST_IMAGE
        
        print(f"[mock] POST /images → created {meta['uuid']} (seed={seed}, name='{name}')")
        return self._send_json(meta, 201)
    
    def do_PUT(self):
        if not self.path.startswith("/images/"):
            return self._send_json(error_response("Not found"), 404)
        
        img_uuid = self._extract_uuid(self.path)
        if img_uuid not in IMAGES:
            return self._send_json(error_response("Image not found"), 404)
        
        try:
            body = self._read_json_body()
        except Exception as e:
            return self._send_json(error_response(f"Invalid JSON: {e}"), 400)
        
        if body is None:
            return self._send_json(error_response("Empty body"), 400)
        
        new_name = (body.get("name") or "").strip()
        if not new_name:
            return self._send_json(error_response("name is required"), 400)
        if len(new_name) > 255:
            return self._send_json(error_response("name too long"), 400)
        
        IMAGES[img_uuid]["name"] = new_name
        print(f"[mock] PUT /images/{img_uuid} → renamed to '{new_name}'")
        return self._send_json(IMAGES[img_uuid])
    
    def do_DELETE(self):
        if not self.path.startswith("/images/"):
            return self._send_json(error_response("Not found"), 404)
        
        img_uuid = self._extract_uuid(self.path)
        if img_uuid not in IMAGES:
            return self._send_json(error_response("Image not found"), 404)
        
        del IMAGES[img_uuid]
        IMAGE_BYTES.pop(img_uuid, None)
        print(f"[mock] DELETE /images/{img_uuid}")
        return self._send_no_content()
    
    def log_message(self, fmt, *args):
        pass

def seed_test_data():
    test_data = [
        (12345, "Neon City"),
        (67890, "Dark Forest"),
        (11111, "Toxic River"),
        (22222, "Rainbow Sky"),
        (33333, "Ocean Deep"),
    ]
    for seed, name in test_data:
        meta = make_image_meta(seed, name)
        IMAGES[meta["uuid"]] = meta
        IMAGE_BYTES[meta["uuid"]] = TEST_IMAGE

if __name__ == "__main__":
    seed_test_data()
    
    PORT = 8080
    print("=" * 60)
    print(f"Mock server for Immersive Images")
    print(f"URL: http://localhost:{PORT}")
    print(f"Loaded {len(IMAGES)} test images")
    print("=" * 60)
    print("Endpoints:")
    print(f"  GET    /healthz")
    print(f"  GET    /images?limit=50&start=<uuid>")
    print(f"  POST   /images")
    print(f"  GET    /images/<uuid>")
    print(f"  PUT    /images/<uuid>")
    print(f"  DELETE /images/<uuid>")
    print(f"  GET    /images/<uuid>/content")
    print("=" * 60)
    print("Press Ctrl+C to stop.")
    print()
    
    try:
        HTTPServer(("0.0.0.0", PORT), Handler).serve_forever()
    except KeyboardInterrupt:
        print("\n[mock] Shutting down.")
        sys.exit(0)