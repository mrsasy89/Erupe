# Erupe wiki import with verified quest-type corrections

This is a staging package for the GitHub wiki at https://github.com/mrsasy89/Erupe/wiki . The wiki has its own Git repository; committing this package to a normal branch does not publish the wiki.

## Contents

- `wiki/`: all 12 pages from Houmgaor/Erupe.wiki snapshot `69ade90d207aacd7cb674b0feb4ca8b266e22e70`, plus a client-mapping evidence page and sidebar.
- `evidence/`: extracted menu descriptors and input identity/check results; no client binaries.
- `Publish-ErupeWiki.ps1`: copies the pages into a fresh checkout of the destination wiki. It preserves unrelated existing pages, backs up overlapping files locally, and uses a normal non-forced Git push. Run without `-Publish` for a local preview.

## Corrections

- Type **54**: **Gear Acquisition Quests**, replacing the misleading upstream alias “G Armor”. Both HR and G entries accept this type, with client eligibility filters.
- Type **9**: **Promotional** in Special Quests; retain Campaign as a server/upstream alias and explain the server's campaign checks.
- Add client-observed types **12, 19, 29, 41, 42, 44, 45**.
- Distinguish server quest types from menu indices and internal list selectors. In particular, internal selector **57** for G Deco Quests is not an SQL quest_type.
- Keep upstream-only associations explicitly unverified; do not invent missing mappings.
- Preserve quest season/time flags, marks, and unrelated wiki pages.

## Publish from Windows

1. If GitHub shows an empty wiki, sign in and create a first page named **Home** at https://github.com/mrsasy89/Erupe/wiki . This initializes the separate wiki repository.
2. Use Git for Windows with your GitHub login and Git author name/email configured.
3. Run `powershell -ExecutionPolicy Bypass -File .\Publish-ErupeWiki.ps1 -Publish` from this package directory.

The script does not modify Erupe code, the database, quest BINs, or the source wiki. Existing pages with the same names are updated to the prepared snapshot. Changes are recoverable through the destination wiki Git history and the local backup.

Original wiki credit remains with Houmgaor/Erupe contributors. Client mappings are version-specific static findings. Type 9 was also tested in game; Type 54 still requires runtime verification. General upstream wiki instructions were copied, not independently tested on this fork.
