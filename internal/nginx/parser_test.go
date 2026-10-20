package nginx_test

import (
	"testing"

	"openresty-admin/internal/nginx"
)

const sample = `server {
    listen 443 ssl http2;
    server_name example.com www.example.com;
    root /var/www/html;
    access_log /var/log/nginx/example.access.log;

    location /api/ {
        proxy_pass http://backend_api;
        access_by_lua_block {
            ngx.say("guarded")
        }
    }

    location /static/ {
        root /var/www/static;
    }
}`

func TestParseSitesExtractsServerLevelDirectives(t *testing.T) {
	sites := nginx.ParseSites(sample, "sites-enabled/example.conf")
	if len(sites) != 1 {
		t.Fatalf("expected 1 site, got %d", len(sites))
	}
	site := sites[0]
	if site.Name != "example.com" {
		t.Fatalf("unexpected name %q", site.Name)
	}
	if len(site.ServerName) != 2 || site.ServerName[1] != "www.example.com" {
		t.Fatalf("unexpected server names %v", site.ServerName)
	}
	if !site.SSL || !site.HTTP2 {
		t.Fatalf("expected ssl+http2, got ssl=%v http2=%v", site.SSL, site.HTTP2)
	}
	if site.Root != "/var/www/html" {
		t.Fatalf("unexpected root %q", site.Root)
	}
	if site.AccessLog != "/var/log/nginx/example.access.log" {
		t.Fatalf("unexpected access log %q", site.AccessLog)
	}
}

func TestParseSitesExtractsLocationsAndLua(t *testing.T) {
	site := nginx.ParseSites(sample, "f.conf")[0]
	if len(site.Locations) != 2 {
		t.Fatalf("expected 2 locations, got %d", len(site.Locations))
	}
	if site.Locations[0].Path != "/api/" {
		t.Fatalf("unexpected location path %q", site.Locations[0].Path)
	}
	if site.Locations[1].Path != "/static/" {
		t.Fatalf("unexpected second location path %q", site.Locations[1].Path)
	}
	if site.Locations[0].ProxyPass != "http://backend_api" {
		t.Fatalf("unexpected proxy_pass %q", site.Locations[0].ProxyPass)
	}
	if len(site.Locations[0].LuaDirectives) != 1 {
		t.Fatalf("expected 1 lua directive, got %d", len(site.Locations[0].LuaDirectives))
	}
	if site.Locations[0].LuaDirectives[0].Type != "access_by_lua_block" {
		t.Fatalf("unexpected directive %q", site.Locations[0].LuaDirectives[0].Type)
	}
}

func TestParseSitesReturnsEmptyWithoutServerBlock(t *testing.T) {
	if sites := nginx.ParseSites("http { }", "nginx.conf"); len(sites) != 0 {
		t.Fatalf("expected no sites, got %d", len(sites))
	}
}
