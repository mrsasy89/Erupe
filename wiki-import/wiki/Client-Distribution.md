# Client Distribution

Erupe emulates the game servers. It does **not** give players a way to reach
them: no launcher, no game files, no server list. This page covers the pieces
an operator hosts *around* Erupe so that a player with an empty folder ends up
in the game, and which `config.json` keys tie them together.

If you only ever connect from your own machine with a client you already have,
you can skip most of this — set `Host`, one `Entrance.Entries` world, and go.

## What a player-facing deployment consists of

| Piece | Who needs it | Configured by |
|-------|--------------|---------------|
| Sign / entrance / channel servers | everyone | `Sign`, `Entrance`, `Channel` — see [Erupe Configuration](Erupe-Configuration) |
| World list | everyone | `Entrance.Entries` |
| REST API (`:8080`) | mhf-outpost players, web sign-up, save export | `API.Enabled`, `API.Port` |
| Patch tree (`mhfdat/{exe,dat}` + manifest) | anyone who does not already have the ~5 GB of game files, and every update you ship | `API.PatchServer`, `PatchServerManifest`, `PatchServerFile` |
| Launcher page + `serverslist.xml` | players on the original `mhf.exe` launcher only | outside Erupe (static files) |

There are two ways for a player to log in, and they need different things
from you. Most servers support both.

```
mhf-outpost ──HTTP──▶ Erupe API :8080  (login, character list, patchServer)
            ──HTTP──▶ patch tree       (game files)
            ──game──▶ sign :53312 → entrance :53310 → channel :54001+

original launcher (mhf.exe + mhl.dll)
            ──HTTP──▶ /launcher/ page + serverslist.xml   (static, hosted by you or a hub)
            ──game──▶ sign :53312  (login; reply carries PatchServerManifest/File)
            ──HTTP──▶ patch tree
            ──game──▶ entrance :53310 → channel :54001+
```

## Path A — players use mhf-outpost

[mhf-outpost](https://github.com/Mogapedia/mhf-outpost) is an open-source
launcher that ships no game data. A player types your API address, signs in,
and the launcher syncs the game files from wherever `API.PatchServer` points.
On your side:

1. Enable the API and tell it where the patch tree is:

   ```json
   "API": {
     "Enabled": true,
     "Port": 8080,
     "PatchServer": "http://frontier.example.com"
   }
   ```

   `PatchServer` is returned verbatim as `patchServer` in `POST /v2/login`.
   It is the origin (scheme + host, optional sub-path) under which
   `/mhf_file.php` and `/mhfdat/` are reachable — see
   [the patch tree](#hosting-the-patch-tree) below.

2. Host the patch tree over plain HTTP.

3. Tell players: *"server address `http://frontier.example.com:8080`"*.
   Nothing else — no hosts file, no `url.lst`, no launcher page.

If you put the API behind a reverse proxy with TLS, players can use `https://`
for the API; the patch tree still has to be reachable over plain `http://`
because the original launcher cannot do TLS, and there is no reason to
maintain two trees.

## Path B — players use the original launcher

The original client boots through `mhf.exe` → `mhl.dll`, a Trident (IE)
web view. `mhl.dll` holds two hard-coded URLs, both plain HTTP:

- `http://<host>/launcher/` — the launcher HTML page
- `http://<host>/serverslist.xml` — the server list shown in the page

Capcom's DLL points at `cog-members.mh-frontier.jp`. Community builds are
hex-patched to point at a hub; the widely distributed one points at the
Mezeporta hub, whose `serverslist.xml` lists the public Erupe servers, and its
`Custom` entry resolves the hostname `erupe.custom` — that is what "edit your
hosts file" means in older guides.

To support these players you either:

- **get listed on a hub's `serverslist.xml`** (ask the hub operator), or
- **host your own launcher page** and distribute an `mhl.dll` patched with
  your host. The page is static HTML/JS from
  [Mezeporta/Servers](https://github.com/Mezeporta/Servers); `serverslist.xml`
  looks like:

  ```xml
  <?xml version="1.0" encoding="shift_jis"?>
  <server_groups>
    <group idx='0' nam='Select a Server' ip=''/>
    <group idx='1' nam='My Server' ip='203.0.113.10' port='53312' svid='mine' oauth_url='debug-www.mhf-g.com/console-link.html'/>
  </server_groups>
  ```

  `ip`/`port` are your sign server. Login goes straight from `mhl.dll` to the
  sign server; the REST API is not involved.

Then set the patch server hosts the sign server hands back on login:

```json
"PatchServerManifest": "frontier.example.com",
"PatchServerFile": "frontier.example.com"
```

Both are **bare hosts** (optionally `host/subdir`), no `http://`. Left
empty, the PC launcher skips the update step, which is fine for players who
already have a complete install; PS3 clients are refused at login
(`SIGN_EABORT`) unless both are set. The two can differ if you serve the
manifest and the files from different machines.

> The launcher page, `serverslist.xml`, the manifest and the files must all be
> served over plain HTTP. If your web server redirects `http://` to
> `https://`, the launcher shows a blank page or fails silently. Keep the
> redirect off for that vhost; you can *also* listen on 443 for browsers.

## Hosting the patch tree

This is the same protocol Capcom's patch server spoke, so it works unchanged
for the original launcher and for mhf-outpost. Requests, as issued by
`mhl.dll` (User-Agent `mhfhttp`):

```
GET http://<PatchServerManifest>/mhf_file.php?key=<n>            manifest, saved locally as MHFUP_00.DAT
GET http://<PatchServerManifest>/mhf_file.php?chk=<n>            expects the text  [mhf Check Message:0]
GET http://<PatchServerFile>/mhfdat/<exe|dat>/<file>?id=<n>      one file
```

`<n>` is a cache-buster; ignore it.

### Layout

```
patch/
├── key.txt          manifest (served as /mhf_file.php?key=…)
├── chk.txt          contains exactly:  [mhf Check Message:0]
└── mhfdat/
    ├── exe/         files that live next to mhf.exe: mhf.exe, mhl.dll, mhfo.dll,
    │                mhfo-hd.dll, gdiplus.dll, ogg.dll, vorbis.dll, url*.lst
    └── dat/         the game's dat/ folder
```

The manifest is one line per file:

```
CRC32,FILETIME_LO,FILETIME_HI,path,size,0
3638A322,B61F3F00,01DA01F4,dat\stage\st277.pac,2867018,0
```

Paths use backslashes and are relative to the game folder (`exe\` entries
are written to the root). The client compares CRC32 and size against what it
has and downloads only the differences; it never deletes files, so removing
something from the tree does not remove it from players' installs.
Zero-byte files must be left out of the manifest.

### Option 1 — let Erupe serve it (simplest)

```json
"API": {
  "Enabled": true,
  "Port": 8080,
  "PatchTree": { "Enabled": true, "Root": "patch" }
}
```

Put the files under `patch/mhfdat/{exe,dat}` (relative to the working
directory, or give an absolute `Root`) and start Erupe. On startup it checks
whether any file is newer than `key.txt` and, if so, regenerates the
manifest in the background (about 4 s for a ZZ tree on a warm disk; the
manifest route answers 503 until it is done). `chk.txt` is not needed. The
tree is then served on the API port:

```
http://<Host>:8080/mhf_file.php?key=1
http://<Host>:8080/mhfdat/exe/mhfo-hd.dll
```

`API.PatchServer` defaults to `http://<Host>:<API.Port>` when `PatchTree` is
enabled, so mhf-outpost needs nothing else. For the original launcher, set
`PatchServerManifest`/`PatchServerFile` to `<Host>:8080`; whether `mhl.dll`
honours a port in that string has not been verified, so if it fails, put the
API behind a reverse proxy on port 80 and point all three settings at the
proxy's address.

Updating the tree is: replace the files, restart Erupe (or run `patchtree gen`
below, no restart needed).

### Option 2 — serve it from a web server

Generate the manifest with either tool from the Erupe repository, run from
the `patch/` directory:

```bash
go build -o patchtree ./cmd/patchtree && ./patchtree gen --root patch   # writes patch/key.txt
# or
cd patch && python3 genMhfKey.py && mv key00.txt key.txt                 # docs/patch-server/genMhfKey.py
printf '[mhf Check Message:0]' > patch/chk.txt
```

Re-run it every time you change a file in the tree (translation update, new
DLL, …). Nothing needs restarting — clients pick up the new manifest on their
next launch.

Any static file server works for `/mhfdat/`; `/mhf_file.php` just needs the
query string mapped to the two text files. A complete nginx example is in
`docs/patch-server/nginx-patch-server.conf`; the core of it:

```nginx
location = /mhf_file.php {
    root /var/www/patch;
    default_type application/octet-stream;
    if ($args ~ "(^|&)key(=|&|$)") { rewrite ^ /key.txt break; }
    if ($args ~ "(^|&)chk(sha)?(=|&|$)") { rewrite ^ /chk.txt break; }
    return 404;
}
location /mhfdat/ {
    root /var/www/patch;
    default_type application/octet-stream;
    autoindex off;
}
```

Sanity check from a shell:

```bash
curl -s 'http://frontier.example.com/mhf_file.php?chk=1'        # [mhf Check Message:0]
curl -s 'http://frontier.example.com/mhf_file.php?key=1' | wc -l  # number of files
curl -sI 'http://frontier.example.com/mhfdat/exe/mhfo-hd.dll'     # 200, Content-Length ≈ 29 MB
```

A ZZ tree is about 5.3 GB and 5 000 files. Put it on a disk that can take
several players doing a first install at once; there is no throttling on the
client side. That is the one reason to prefer option 2 on a busy server:
nginx serving static files costs less than the Go process doing it, and it
keeps download traffic off the API port.

### What goes in the tree

Whatever you want players to end up with: the client version matching your
`ClientMode` (see [Client Versions](Client-Versions)), your patched
`mhl.dll` if you host a launcher page, and any translation or modified data
files. Erupe does not ship any of this and the project makes no statement
about what you may serve — that is between you and the rights holders.

## Choosing a client version

`ClientMode` in `config.json` must match the client build in the patch tree
(and on players' disks). Mixing them produces silent failures at login or in
quest lists. The [Client Versions](Client-Versions) table maps each release to
its `ClientMode` ID; `GET /v2/server/info` reports the running value so a
launcher can check before syncing.

## Ports to open

| Port | Protocol | Purpose |
|------|----------|---------|
| 53312 | TCP | Sign server |
| 53310 | TCP | Entrance server |
| 54001+ | TCP | One per channel in `Entrance.Entries[].Channels` |
| 8080 | TCP | REST API (mhf-outpost, sign-up pages); also the patch tree with `API.PatchTree` |
| 80 | TCP | Patch tree and launcher page when served by a web server (plain HTTP, required) |

See [Server Operations → Firewall](Server-Operations#firewall-configuration)
for `ufw` examples.

## Checklist

- [ ] `Host` is the public IP/hostname; `Entrance.Entries[].IP` too
- [ ] `ClientMode` matches the client build you distribute
- [ ] `API.Enabled: true`; `API.PatchTree.Enabled: true` or `API.PatchServer` set (mhf-outpost)
- [ ] `PatchServerManifest` / `PatchServerFile` set (original launcher)
- [ ] `key.txt` regenerated after the last change to `mhfdat/`
- [ ] `curl 'http://<host>/mhf_file.php?chk=1'` returns `[mhf Check Message:0]` over **plain HTTP**
- [ ] no `http → https` redirect on that vhost
- [ ] players told either the API address (mhf-outpost) or where to get the launcher (original)
