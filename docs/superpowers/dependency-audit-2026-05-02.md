# Dependency Audit — 2026-05-02

Baseline captures for the [dependency refresh plan](plans/2026-05-02-dependency-refresh.md).

## go list -u -m all

```
github.com/guarzo/canifly
cloud.google.com/go/compute/metadata v0.3.0 [v0.9.0]
github.com/davecgh/go-spew v1.1.1
github.com/felixge/httpsnoop v1.0.4
github.com/google/gofuzz v1.2.0
github.com/google/uuid v1.6.0
github.com/gorilla/handlers v1.5.2
github.com/gorilla/mux v1.8.1
github.com/gorilla/securecookie v1.1.2
github.com/gorilla/sessions v1.4.0
github.com/gorilla/websocket v1.5.3
github.com/pmezard/go-difflib v1.0.0
github.com/sirupsen/logrus v1.9.3 [v1.9.4]
github.com/stretchr/objx v0.5.2 [v0.5.3]
github.com/stretchr/testify v1.10.0 [v1.11.1]
golang.org/x/oauth2 v0.30.0 [v0.36.0]
golang.org/x/sys v0.35.0 [v0.43.0]
gopkg.in/check.v1 v0.0.0-20161208181325-20d25e280405 [v1.0.0-20201130134442-10cb98267c6c]
gopkg.in/yaml.v3 v3.0.1
```

## govulncheck ./...

```
=== Symbol Results ===

Vulnerability #1: GO-2026-4947
    Unexpected work during chain building in crypto/x509
  More info: https://pkg.go.dev/vuln/GO-2026-4947
  Standard library
    Found in: crypto/x509@go1.25.7
    Fixed in: crypto/x509@go1.25.9
    Example traces found:
      #1: internal/persist/crypto.go:123:26: persist.EncryptString calls io.ReadFull, which eventually calls x509.Certificate.Verify

Vulnerability #2: GO-2026-4946
    Inefficient policy validation in crypto/x509
  More info: https://pkg.go.dev/vuln/GO-2026-4946
  Standard library
    Found in: crypto/x509@go1.25.7
    Fixed in: crypto/x509@go1.25.9
    Example traces found:
      #1: internal/persist/crypto.go:123:26: persist.EncryptString calls io.ReadFull, which eventually calls x509.Certificate.Verify

Vulnerability #3: GO-2026-4870
    Unauthenticated TLS 1.3 KeyUpdate record can cause persistent connection
    retention and DoS in crypto/tls
  More info: https://pkg.go.dev/vuln/GO-2026-4870
  Standard library
    Found in: crypto/tls@go1.25.7
    Fixed in: crypto/tls@go1.25.9
    Example traces found:
      #1: internal/cmd/start.go:103:19: cmd.runServer calls http.Server.Serve, which eventually calls tls.Conn.HandshakeContext
      #2: internal/persist/crypto.go:123:26: persist.EncryptString calls io.ReadFull, which eventually calls tls.Conn.Read
      #3: internal/services/account/auth_client.go:99:2: account.authClient.RefreshToken calls http.http2transportResponseBody.Close, which eventually calls tls.Conn.Write
      #4: internal/services/account/auth_client.go:94:26: account.authClient.RefreshToken calls http.Client.Do, which eventually calls tls.Dialer.DialContext

Vulnerability #4: GO-2026-4602
    FileInfo can escape from a Root in os
  More info: https://pkg.go.dev/vuln/GO-2026-4602
  Standard library
    Found in: os@go1.25.7
    Fixed in: os@go1.25.8
    Example traces found:
      #1: internal/cmd/start.go:103:19: cmd.runServer calls http.Server.Serve, which eventually calls os.File.ReadDir
      #2: internal/cmd/start.go:103:19: cmd.runServer calls http.Server.Serve, which eventually calls os.File.Readdir
      #3: internal/persist/eve/eve_profiles_store.go:164:28: eve.EveProfilesStore.GetSubDirectories calls os.ReadDir

Vulnerability #5: GO-2026-4601
    Incorrect parsing of IPv6 host literals in net/url
  More info: https://pkg.go.dev/vuln/GO-2026-4601
  Standard library
    Found in: net/url@go1.25.7
    Fixed in: net/url@go1.25.8
    Example traces found:
      #1: internal/services/account/auth_client.go:85:29: account.authClient.RefreshToken calls http.NewRequest, which eventually calls url.Parse
      #2: internal/cmd/start.go:103:19: cmd.runServer calls http.Server.Serve, which eventually calls url.ParseRequestURI
      #3: internal/services/account/auth_client.go:94:26: account.authClient.RefreshToken calls http.Client.Do, which eventually calls url.URL.Parse

Your code is affected by 5 vulnerabilities from the Go standard library.
This scan also found 2 vulnerabilities in packages you import and 2
vulnerabilities in modules you require, but your code doesn't appear to call
these vulnerabilities.
Use '-show verbose' for more details.
```

## npm audit (root)

```
# npm audit report

@isaacs/brace-expansion  5.0.0
Severity: high
@isaacs/brace-expansion has Uncontrolled Resource Consumption - https://github.com/advisories/GHSA-7h2j-956f-4vf2
fix available via `npm audit fix`
node_modules/@isaacs/brace-expansion

@tootallnate/once  <3.0.1
@tootallnate/once vulnerable to Incorrect Control Flow Scoping - https://github.com/advisories/GHSA-vpq2-c234-7xj6
fix available via `npm audit fix --force`
Will install electron-builder@26.8.1, which is a breaking change
node_modules/@tootallnate/once
  http-proxy-agent  4.0.1 - 5.0.0
  Depends on vulnerable versions of @tootallnate/once
  node_modules/make-fetch-happen/node_modules/http-proxy-agent
    make-fetch-happen  7.1.1 - 14.0.0
    Depends on vulnerable versions of cacache
    Depends on vulnerable versions of http-proxy-agent
    node_modules/make-fetch-happen
      node-gyp  <=10.3.1
      Depends on vulnerable versions of make-fetch-happen
      Depends on vulnerable versions of tar
      node_modules/node-gyp
        @electron/rebuild  3.2.10 - 4.0.2
        Depends on vulnerable versions of node-gyp
        Depends on vulnerable versions of tar
        node_modules/@electron/rebuild
          app-builder-lib  23.0.7 - 26.5.0
          Depends on vulnerable versions of @electron/rebuild
          Depends on vulnerable versions of dmg-builder
          Depends on vulnerable versions of electron-builder-squirrel-windows
          Depends on vulnerable versions of tar
          node_modules/app-builder-lib
            dmg-builder  23.0.7 - 26.5.0
            Depends on vulnerable versions of app-builder-lib
            node_modules/dmg-builder
              electron-builder  19.25.0 || 23.0.7 - 26.5.0
              Depends on vulnerable versions of app-builder-lib
              Depends on vulnerable versions of dmg-builder
              node_modules/electron-builder
            electron-builder-squirrel-windows  23.0.7 - 26.5.0
            Depends on vulnerable versions of app-builder-lib
            node_modules/electron-builder-squirrel-windows

@xmldom/xmldom  <=0.8.12
Severity: high
xmldom: Uncontrolled recursion in XML serialization leads to DoS - https://github.com/advisories/GHSA-2v35-w6hq-6mfw
xmldom has XML injection through unvalidated DocumentType serialization - https://github.com/advisories/GHSA-f6ww-3ggp-fr8h
xmldom has XML node injection through unvalidated processing instruction serialization - https://github.com/advisories/GHSA-x6wf-f3px-wcqx
xmldom has XML node injection through unvalidated comment serialization - https://github.com/advisories/GHSA-j759-j44w-7fr8
xmldom: XML injection via unsafe CDATA serialization allows attacker-controlled markup insertion - https://github.com/advisories/GHSA-wh4c-j3r5-mjhp
fix available via `npm audit fix`
node_modules/@xmldom/xmldom

ajv  <6.14.0
Severity: moderate
ajv has ReDoS when using `$data` option - https://github.com/advisories/GHSA-2g4f-4pwh-qvx6
fix available via `npm audit fix`
node_modules/ajv

axios  1.0.0 - 1.14.0
Severity: high
Axios is vulnerable to DoS attack through lack of data size check - https://github.com/advisories/GHSA-4hjh-wcwx-xvwj
Axios is Vulnerable to Denial of Service via __proto__ Key in mergeConfig - https://github.com/advisories/GHSA-43fc-jf86-j433
Axios has a NO_PROXY Hostname Normalization Bypass that Leads to SSRF - https://github.com/advisories/GHSA-3p68-rc4w-qgx5
Axios has Unrestricted Cloud Metadata Exfiltration via Header Injection Chain - https://github.com/advisories/GHSA-fvcv-3m26-pcqx
fix available via `npm audit fix`
node_modules/axios

brace-expansion  <1.1.13 || >=2.0.0 <2.0.3
Severity: moderate
brace-expansion: Zero-step sequence causes process hang and memory exhaustion - https://github.com/advisories/GHSA-f886-m6hf-6m8v
brace-expansion: Zero-step sequence causes process hang and memory exhaustion - https://github.com/advisories/GHSA-f886-m6hf-6m8v
fix available via `npm audit fix`
node_modules/@electron/universal/node_modules/brace-expansion
node_modules/brace-expansion
node_modules/cacache/node_modules/brace-expansion
node_modules/config-file-ts/node_modules/brace-expansion
node_modules/filelist/node_modules/brace-expansion
node_modules/readdir-glob/node_modules/brace-expansion

electron  <=39.8.4
Severity: high
Electron has ASAR Integrity Bypass via resource modification - https://github.com/advisories/GHSA-vmqv-hx8q-j7mg
Electron: AppleScript injection in app.moveToApplicationsFolder on macOS - https://github.com/advisories/GHSA-5rqw-r77c-jp79
Electron: Service worker can spoof executeJavaScript IPC replies - https://github.com/advisories/GHSA-xj5x-m3f3-5x3h
Electron: Incorrect origin passed to permission request handler for iframe requests - https://github.com/advisories/GHSA-r5p7-gp4j-qhrx
Electron: Out-of-bounds read in second-instance IPC on macOS and Linux - https://github.com/advisories/GHSA-3c8v-cfp5-9885
Electron: nodeIntegrationInWorker not correctly scoped in shared renderer processes - https://github.com/advisories/GHSA-xwr5-m59h-vwqr
Electron: Use-after-free in offscreen child window paint callback - https://github.com/advisories/GHSA-532v-xpq5-8h95
Electron: Registry key path injection in app.setAsDefaultProtocolClient on Windows - https://github.com/advisories/GHSA-mwmh-mq4g-g6gr
Electron: Use-after-free in download save dialog callback - https://github.com/advisories/GHSA-9w97-2464-8783
Electron: Use-after-free in WebContents fullscreen, pointer-lock, and keyboard-lock permission callbacks - https://github.com/advisories/GHSA-8337-3p73-46f4
Electron: Use-after-free in PowerMonitor on Windows and macOS - https://github.com/advisories/GHSA-jjp3-mq3x-295m
Electron: Unquoted executable path in app.setLoginItemSettings on Windows - https://github.com/advisories/GHSA-jfqx-fxh3-c62j
Electron: HTTP Response Header Injection in custom protocol handlers and webRequest - https://github.com/advisories/GHSA-4p4r-m79c-wq3v
Electron: USB device selection not validated against filtered device list - https://github.com/advisories/GHSA-9899-m83m-qhpj
Electron: Crash in clipboard.readImage() on malformed clipboard image data - https://github.com/advisories/GHSA-f37v-82c4-4x64
Electron: Named window.open targets not scoped to the opener's browsing context - https://github.com/advisories/GHSA-f3pv-wv63-48x8
Electron: Renderer command-line switch injection via undocumented commandLineSwitches webPreference - https://github.com/advisories/GHSA-9wfr-w7mm-pc7f
fix available via `npm audit fix --force`
Will install electron@41.5.0, which is a breaking change
node_modules/electron

flatted  <=3.4.1
Severity: high
flatted vulnerable to unbounded recursion DoS in parse() revive phase - https://github.com/advisories/GHSA-25h7-pfq9-p65f
Prototype Pollution via parse() in NodeJS flatted - https://github.com/advisories/GHSA-rf6f-7fwh-wjgh
fix available via `npm audit fix`
node_modules/flatted

follow-redirects  <=1.15.11
Severity: moderate
follow-redirects leaks Custom Authentication Headers to Cross-Domain Redirect Targets - https://github.com/advisories/GHSA-r4q5-vmmm-2653
fix available via `npm audit fix`
node_modules/follow-redirects

glob  10.2.0 - 10.4.5
Severity: high
glob CLI: Command injection via -c/--cmd executes matches with shell:true - https://github.com/advisories/GHSA-5j98-mcp5-4vw2
fix available via `npm audit fix`
node_modules/config-file-ts/node_modules/glob

handlebars  4.0.0 - 4.7.8
Severity: critical
Handlebars.js has JavaScript Injection via AST Type Confusion by tampering @partial-block - https://github.com/advisories/GHSA-3mfm-83xf-c92r
Handlebars.js has JavaScript Injection via AST Type Confusion - https://github.com/advisories/GHSA-2w6w-674q-4c4q
Handlebars.js has Prototype Pollution Leading to XSS through Partial Template Injection - https://github.com/advisories/GHSA-2qvq-rjwj-gvw9
Handlebars.js has a Prototype Method Access Control Gap via Missing __lookupSetter__ Blocklist Entry - https://github.com/advisories/GHSA-7rx3-28cr-v5wh
Handlebars.js has a Property Access Validation Bypass in container.lookup - https://github.com/advisories/GHSA-442j-39wm-28r2
Handlebars.js has JavaScript Injection via AST Type Confusion when passing an object as dynamic partial - https://github.com/advisories/GHSA-xhpv-hc6g-r9c6
Handlebars.js has Denial of Service via Malformed Decorator Syntax in Template Compilation - https://github.com/advisories/GHSA-9cx6-37pm-9jff
Handlebars.js has JavaScript Injection in CLI Precompiler via Unescaped Names and Options - https://github.com/advisories/GHSA-xjpj-3mr7-gcpf
fix available via `npm audit fix`
node_modules/handlebars

js-yaml  4.0.0 - 4.1.0
Severity: moderate
js-yaml has prototype pollution in merge (<<) - https://github.com/advisories/GHSA-mh29-5h37-fv8m
fix available via `npm audit fix`
node_modules/js-yaml

lodash  <=4.17.23
Severity: high
Lodash has Prototype Pollution Vulnerability in `_.unset` and `_.omit` functions - https://github.com/advisories/GHSA-xxjr-mmjv-4gpg
lodash vulnerable to Code Injection via `_.template` imports key names - https://github.com/advisories/GHSA-r5fr-rjxr-66jc
lodash vulnerable to Prototype Pollution via array path bypass in `_.unset` and `_.omit` - https://github.com/advisories/GHSA-f23m-r3pf-42rh
fix available via `npm audit fix`
node_modules/lodash

minimatch  <=3.1.3 || 5.0.0 - 5.1.7 || 9.0.0 - 9.0.6 || 10.0.0 - 10.2.2
Severity: high
minimatch has a ReDoS via repeated wildcards with non-matching literal in pattern - https://github.com/advisories/GHSA-3ppc-4f35-3m26
minimatch has a ReDoS via repeated wildcards with non-matching literal in pattern - https://github.com/advisories/GHSA-3ppc-4f35-3m26
minimatch has a ReDoS via repeated wildcards with non-matching literal in pattern - https://github.com/advisories/GHSA-3ppc-4f35-3m26
minimatch has a ReDoS via repeated wildcards with non-matching literal in pattern - https://github.com/advisories/GHSA-3ppc-4f35-3m26
minimatch has ReDoS: matchOne() combinatorial backtracking via multiple non-adjacent GLOBSTAR segments - https://github.com/advisories/GHSA-7r86-cg39-jmmj
minimatch has ReDoS: matchOne() combinatorial backtracking via multiple non-adjacent GLOBSTAR segments - https://github.com/advisories/GHSA-7r86-cg39-jmmj
minimatch has ReDoS: matchOne() combinatorial backtracking via multiple non-adjacent GLOBSTAR segments - https://github.com/advisories/GHSA-7r86-cg39-jmmj
minimatch has ReDoS: matchOne() combinatorial backtracking via multiple non-adjacent GLOBSTAR segments - https://github.com/advisories/GHSA-7r86-cg39-jmmj
minimatch ReDoS: nested *() extglobs generate catastrophically backtracking regular expressions - https://github.com/advisories/GHSA-23c5-xmqv-rm74
minimatch ReDoS: nested *() extglobs generate catastrophically backtracking regular expressions - https://github.com/advisories/GHSA-23c5-xmqv-rm74
minimatch ReDoS: nested *() extglobs generate catastrophically backtracking regular expressions - https://github.com/advisories/GHSA-23c5-xmqv-rm74
minimatch ReDoS: nested *() extglobs generate catastrophically backtracking regular expressions - https://github.com/advisories/GHSA-23c5-xmqv-rm74
fix available via `npm audit fix`
node_modules/@electron/asar/node_modules/minimatch
node_modules/@electron/universal/node_modules/minimatch
node_modules/@eslint/config-array/node_modules/minimatch
node_modules/@eslint/eslintrc/node_modules/minimatch
node_modules/cacache/node_modules/minimatch
node_modules/config-file-ts/node_modules/minimatch
node_modules/dir-compare/node_modules/minimatch
node_modules/dotgitignore/node_modules/minimatch
node_modules/eslint/node_modules/minimatch
node_modules/filelist/node_modules/minimatch
node_modules/glob/node_modules/minimatch
node_modules/minimatch
node_modules/readdir-glob/node_modules/minimatch

picomatch  4.0.0 - 4.0.3
Severity: high
Picomatch: Method Injection in POSIX Character Classes causes incorrect Glob Matching - https://github.com/advisories/GHSA-3v7f-55p6-f55p
Picomatch has a ReDoS vulnerability via extglob quantifiers - https://github.com/advisories/GHSA-c2c7-rcm5-vvqj
fix available via `npm audit fix`
node_modules/picomatch

postcss  <8.5.10
Severity: moderate
PostCSS has XSS via Unescaped </style> in its CSS Stringify Output - https://github.com/advisories/GHSA-qx2v-qp2m-jg93
fix available via `npm audit fix`
node_modules/postcss

rollup  4.0.0 - 4.58.0
Severity: high
Rollup 4 has Arbitrary File Write via Path Traversal - https://github.com/advisories/GHSA-mw96-cpmx-2vgc
fix available via `npm audit fix`
node_modules/rollup

tar  <=7.5.10
Severity: high
node-tar Vulnerable to Arbitrary File Creation/Overwrite via Hardlink Path Traversal - https://github.com/advisories/GHSA-34x7-hfp2-rc4v
node-tar is Vulnerable to Arbitrary File Overwrite and Symlink Poisoning via Insufficient Path Sanitization - https://github.com/advisories/GHSA-8qq5-rm4j-mr97
Arbitrary File Read/Write via Hardlink Target Escape Through Symlink Chain in node-tar Extraction - https://github.com/advisories/GHSA-83g3-92jg-28cx
tar has Hardlink Path Traversal via Drive-Relative Linkpath - https://github.com/advisories/GHSA-qffp-2rhf-9h96
node-tar Symlink Path Traversal via Drive-Relative Linkpath - https://github.com/advisories/GHSA-9ppj-qmqm-q256
Race Condition in node-tar Path Reservations via Unicode Ligature Collisions on macOS APFS - https://github.com/advisories/GHSA-r6q2-hw4h-h46w
fix available via `npm audit fix --force`
Will install electron-builder@26.8.1, which is a breaking change
node_modules/@electron/rebuild/node_modules/tar
node_modules/app-builder-lib/node_modules/tar
node_modules/cacache/node_modules/tar
node_modules/node-gyp/node_modules/tar
node_modules/tar
  cacache  14.0.0 - 18.0.4
  Depends on vulnerable versions of tar
  node_modules/cacache

vite  <=6.4.1
Severity: high
Vite middleware may serve files starting with the same name with the public directory - https://github.com/advisories/GHSA-g4jq-h2w9-997c
Vite's `server.fs` settings were not applied to HTML files - https://github.com/advisories/GHSA-jqfw-vq24-v9c3
vite allows server.fs.deny bypass via backslash on Windows - https://github.com/advisories/GHSA-93m4-6634-74q7
Vite Vulnerable to Path Traversal in Optimized Deps `.map` Handling - https://github.com/advisories/GHSA-4w7w-66w2-5vf9
Vite Vulnerable to Arbitrary File Read via Vite Dev Server WebSocket - https://github.com/advisories/GHSA-p9ff-h696-f583
fix available via `npm audit fix`
node_modules/vite

28 vulnerabilities (2 low, 5 moderate, 20 high, 1 critical)

To address issues that do not require attention, run:
  npm audit fix

To address all issues (including breaking changes), run:
  npm audit fix --force
```

## npm audit (renderer)

```
# npm audit report

ajv  <6.14.0
Severity: moderate
ajv has ReDoS when using `$data` option - https://github.com/advisories/GHSA-2g4f-4pwh-qvx6
fix available via `npm audit fix`
node_modules/ajv

axios  1.0.0 - 1.14.0
Severity: high
Axios is vulnerable to DoS attack through lack of data size check - https://github.com/advisories/GHSA-4hjh-wcwx-xvwj
Axios is Vulnerable to Denial of Service via __proto__ Key in mergeConfig - https://github.com/advisories/GHSA-43fc-jf86-j433
Axios has a NO_PROXY Hostname Normalization Bypass that Leads to SSRF - https://github.com/advisories/GHSA-3p68-rc4w-qgx5
Axios has Unrestricted Cloud Metadata Exfiltration via Header Injection Chain - https://github.com/advisories/GHSA-fvcv-3m26-pcqx
fix available via `npm audit fix`
node_modules/axios

brace-expansion  <1.1.13 || >=2.0.0 <2.0.3
Severity: moderate
brace-expansion: Zero-step sequence causes process hang and memory exhaustion - https://github.com/advisories/GHSA-f886-m6hf-6m8v
brace-expansion: Zero-step sequence causes process hang and memory exhaustion - https://github.com/advisories/GHSA-f886-m6hf-6m8v
fix available via `npm audit fix`
node_modules/brace-expansion
node_modules/glob/node_modules/brace-expansion

flatted  <=3.4.1
Severity: high
flatted vulnerable to unbounded recursion DoS in parse() revive phase - https://github.com/advisories/GHSA-25h7-pfq9-p65f
Prototype Pollution via parse() in NodeJS flatted - https://github.com/advisories/GHSA-rf6f-7fwh-wjgh
fix available via `npm audit fix`
node_modules/flatted

follow-redirects  <=1.15.11
Severity: moderate
follow-redirects leaks Custom Authentication Headers to Cross-Domain Redirect Targets - https://github.com/advisories/GHSA-r4q5-vmmm-2653
fix available via `npm audit fix`
node_modules/follow-redirects

glob  10.2.0 - 10.4.5
Severity: high
glob CLI: Command injection via -c/--cmd executes matches with shell:true - https://github.com/advisories/GHSA-5j98-mcp5-4vw2
fix available via `npm audit fix`
node_modules/glob

js-yaml  4.0.0 - 4.1.0
Severity: moderate
js-yaml has prototype pollution in merge (<<) - https://github.com/advisories/GHSA-mh29-5h37-fv8m
fix available via `npm audit fix`
node_modules/js-yaml

lodash  <=4.17.23
Severity: high
Lodash has Prototype Pollution Vulnerability in `_.unset` and `_.omit` functions - https://github.com/advisories/GHSA-xxjr-mmjv-4gpg
lodash vulnerable to Code Injection via `_.template` imports key names - https://github.com/advisories/GHSA-r5fr-rjxr-66jc
lodash vulnerable to Prototype Pollution via array path bypass in `_.unset` and `_.omit` - https://github.com/advisories/GHSA-f23m-r3pf-42rh
fix available via `npm audit fix`
node_modules/lodash

minimatch  <=3.1.3 || 9.0.0 - 9.0.6
Severity: high
minimatch has a ReDoS via repeated wildcards with non-matching literal in pattern - https://github.com/advisories/GHSA-3ppc-4f35-3m26
minimatch has a ReDoS via repeated wildcards with non-matching literal in pattern - https://github.com/advisories/GHSA-3ppc-4f35-3m26
minimatch has ReDoS: matchOne() combinatorial backtracking via multiple non-adjacent GLOBSTAR segments - https://github.com/advisories/GHSA-7r86-cg39-jmmj
minimatch has ReDoS: matchOne() combinatorial backtracking via multiple non-adjacent GLOBSTAR segments - https://github.com/advisories/GHSA-7r86-cg39-jmmj
minimatch ReDoS: nested *() extglobs generate catastrophically backtracking regular expressions - https://github.com/advisories/GHSA-23c5-xmqv-rm74
minimatch ReDoS: nested *() extglobs generate catastrophically backtracking regular expressions - https://github.com/advisories/GHSA-23c5-xmqv-rm74
fix available via `npm audit fix`
node_modules/glob/node_modules/minimatch
node_modules/minimatch

picomatch  <=2.3.1 || 4.0.0 - 4.0.3
Severity: high
Picomatch: Method Injection in POSIX Character Classes causes incorrect Glob Matching - https://github.com/advisories/GHSA-3v7f-55p6-f55p
Picomatch: Method Injection in POSIX Character Classes causes incorrect Glob Matching - https://github.com/advisories/GHSA-3v7f-55p6-f55p
Picomatch has a ReDoS vulnerability via extglob quantifiers - https://github.com/advisories/GHSA-c2c7-rcm5-vvqj
Picomatch has a ReDoS vulnerability via extglob quantifiers - https://github.com/advisories/GHSA-c2c7-rcm5-vvqj
fix available via `npm audit fix`
node_modules/picomatch
node_modules/tinyglobby/node_modules/picomatch
node_modules/vite/node_modules/picomatch
node_modules/vitest/node_modules/picomatch

postcss  <8.5.10
Severity: moderate
PostCSS has XSS via Unescaped </style> in its CSS Stringify Output - https://github.com/advisories/GHSA-qx2v-qp2m-jg93
fix available via `npm audit fix`
node_modules/postcss

react-router  7.0.0 - 7.12.0-pre.0
Severity: high
React Router has CSRF issue in Action/Server Action Request Processing - https://github.com/advisories/GHSA-h5cw-625j-3rxh
React Router vulnerable to XSS via Open Redirects - https://github.com/advisories/GHSA-2w69-qvjg-hvjx
React Router SSR XSS in ScrollRestoration - https://github.com/advisories/GHSA-8v8x-cx79-35w7
React Router has unexpected external redirect via untrusted paths - https://github.com/advisories/GHSA-9jcx-v3wj-wh4m
React Router has XSS Vulnerability - https://github.com/advisories/GHSA-3cgp-3xvw-98x8
fix available via `npm audit fix`
node_modules/react-router
  react-router-dom  7.0.0-pre.0 - 7.11.0
  Depends on vulnerable versions of react-router
  node_modules/react-router-dom

rollup  4.0.0 - 4.58.0
Severity: high
Rollup 4 has Arbitrary File Write via Path Traversal - https://github.com/advisories/GHSA-mw96-cpmx-2vgc
fix available via `npm audit fix`
node_modules/rollup

vite  <=6.4.1
Severity: high
Vite middleware may serve files starting with the same name with the public directory - https://github.com/advisories/GHSA-g4jq-h2w9-997c
Vite's `server.fs` settings were not applied to HTML files - https://github.com/advisories/GHSA-jqfw-vq24-v9c3
vite allows server.fs.deny bypass via backslash on Windows - https://github.com/advisories/GHSA-93m4-6634-74q7
Vite Vulnerable to Path Traversal in Optimized Deps `.map` Handling - https://github.com/advisories/GHSA-4w7w-66w2-5vf9
Vite Vulnerable to Arbitrary File Read via Vite Dev Server WebSocket - https://github.com/advisories/GHSA-p9ff-h696-f583
fix available via `npm audit fix`
node_modules/vite

yaml  1.0.0 - 1.10.2 || 2.0.0 - 2.8.2
Severity: moderate
yaml is vulnerable to Stack Overflow via deeply nested YAML collections - https://github.com/advisories/GHSA-48c2-rrv3-qjmp
yaml is vulnerable to Stack Overflow via deeply nested YAML collections - https://github.com/advisories/GHSA-48c2-rrv3-qjmp
fix available via `npm audit fix`
node_modules/cosmiconfig/node_modules/yaml
node_modules/yaml

16 vulnerabilities (7 moderate, 9 high)

To address all issues, run:
  npm audit fix
```

## npm outdated (root)

```
Package           Current   Wanted   Latest  Location  Depended by
@electron/remote  MISSING    2.1.3    2.1.3  -         dependency-refresh
axios             MISSING   1.15.2   1.15.2  -         dependency-refresh
clsx              MISSING    2.1.1    2.1.1  -         dependency-refresh
dotenv            MISSING   16.6.1   17.4.2  -         dependency-refresh
framer-motion     MISSING  12.38.0  12.38.0  -         dependency-refresh
tailwind-merge    MISSING    3.5.0    3.5.0  -         dependency-refresh
tar               MISSING   7.5.13   7.5.13  -         dependency-refresh
vite              MISSING    6.4.2   8.0.10  -         dependency-refresh
zustand           MISSING   5.0.12   5.0.12  -         dependency-refresh
```

## npm outdated (renderer)

```
Package               Current                                Wanted        Latest  Location  Depended by
@emotion/react        MISSING                               11.14.0       11.14.0  -         renderer
@emotion/styled       MISSING                               11.14.1       11.14.1  -         renderer
@mui/icons-material   MISSING                                 6.5.0         9.0.0  -         renderer
@mui/lab              MISSING  6.0.0-dev.20240529-082515-213b5e33ab  9.0.0-beta.2  -         renderer
@mui/material         MISSING                                 6.5.0         9.0.0  -         renderer
@mui/styles           MISSING                                 6.5.0         6.5.0  -         renderer
electron-is-dev       MISSING                                 3.0.1         3.0.1  -         renderer
framer-motion         MISSING                               12.38.0       12.38.0  -         renderer
next-themes           MISSING                                 0.4.6         0.4.6  -         renderer
prop-types            MISSING                                15.8.1        15.8.1  -         renderer
react                 MISSING                                18.3.1        19.2.5  -         renderer
react-dom             MISSING                                18.3.1        19.2.5  -         renderer
react-error-boundary  MISSING                                 6.1.1         6.1.1  -         renderer
react-icons           MISSING                                 5.6.0         5.6.0  -         renderer
react-router-dom      MISSING                                7.14.2        7.14.2  -         renderer
react-toastify        MISSING                                10.0.6        11.1.0  -         renderer
toastr                MISSING                                 2.1.4         2.1.4  -         renderer
```

## Notes

- **Modules to bump in PR 3:** All five govulncheck findings (GO-2026-4947, GO-2026-4946, GO-2026-4870, GO-2026-4602, GO-2026-4601) are in the Go standard library — fixed by upgrading the Go toolchain to 1.25.9 in PR 2. Two additional vulnerabilities are reported in imported packages and two in required modules (not called from our code); will inspect with `govulncheck -show verbose` in PR 3 and bump root causes there. The npm audit findings collapse to two breaking-change roots: `electron` (root deps) and `electron-builder` (root deps), each pulling many transitives — addressed in PR 5 (electron) and PR 8 (transitive cleanup).
- **Electron major bump path:** currently `^32.2.0`; latest on npm is `41.5.0`. PR 5 walks 32 → 41 (review breaking changes per major).
- **React 19 readiness:** `react` and `react-dom` are at 18.3.1; latest is 19.2.5. PR 6 upgrades to 19.x.
- **MUI v7 readiness:** `@mui/material` (and `@mui/icons-material`, `@mui/lab`) are at 6.5.0; latest is 9.0.0. PR 7 covers the v7+ migration (note: the plan named PR 7 "MUI v7" but current latest is 9.0.0 — will reassess target version when executing PR 7).
