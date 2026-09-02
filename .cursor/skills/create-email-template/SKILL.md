---
name: create-email-template
description: >-
  Creates or updates React Email templates in emails/emails/go_boilerplate from Figma
  Protocol style assets (SaaS Email Templates community file). Use when adding a
  new transactional email, switching email theme, or implementing a Figma email
  frame as code. Requires figma-design-to-code — never hand-draw icons or skip assets.
---

# Create Email Template (Figma Protocol → Go Boilerplate emails)

Add or update transactional email templates in **go.boilerplate** using the **Protocol** style from [SaaS Email Templates (Figma Community)](https://www.figma.com/community/file/1626680546446620209/saas-email-templates).

Brand name in copy: **Go Boilerplate** (never `go.boilerplate`).

## Prerequisites

- Figma MCP authenticated
- Load **figma-design-to-code** before any Figma MCP call
- Read `docs/DESIGN_DOC.md` email section
- `make emails-install` completed once

## Figma reference (Protocol)

| Item | Value |
|------|-------|
| File key | `vkeMjKNcfWdUvmfHdUfL8H` |
| Style section | node `2:7184` |
| Verify / activation | node `2:7186` |
| Password reset | node `2:7395` |
| Design URL | [Protocol selection](https://www.figma.com/design/vkeMjKNcfWdUvmfHdUfL8H/SaaS-Email-Templates--Community-?node-id=2-7184) |

## Workflow (mandatory order)

### 1. get_design_context

Call Figma MCP `get_design_context` on the **specific template frame**.

- Pass `skillNames: figma-design-to-code`
- Extract tokens into `emails/emails/go_boilerplate/theme.ts`
- Treat returned React+Tailwind as **reference only**

### 2. download_assets

Call Figma MCP `download_assets` on the same node.

- Never hand-write SVG paths
- Never commit expiring `figma.com/api/mcp/asset/...` URLs
- Download bytes to `emails/static/go_boilerplate/` using **underscore** file names
- Update `emails/static/go_boilerplate/assets.manifest.json`

### 3. Implement template

Location: `emails/emails/go_boilerplate/<slug>.tsx` (underscores, not dashes)

Reuse: `theme.ts`, `assets.ts`, `go_boilerplate_fonts.tsx`, `go_boilerplate_header.tsx`, `go_boilerplate_footer.tsx`, `go_boilerplate_otp_block.tsx`

```tsx
import { goBoilerplateAsset } from "./assets";
<Img src={goBoilerplateAsset("logo_mark.svg")} width={23} height={23} />
```

### 4. Register renderer

Add the template to `emails/render.ts` with an underscore slug (`verify_email`, `password_reset`).

If Go sends it, update `src/mail/message.go` and `src/mail/service.go`.

### 5. Verify

```bash
make emails-render TEMPLATE=verify_email PROPS='{"code":"123456","name":"Ada","companyName":"Go Boilerplate"}'
make emails-dev
```

Canonical assets live in `emails/static/go_boilerplate/`. Fiber serves them at `/static`. React Email preview serves `--dir`/static (`emails/emails/static`); `make emails-dev` creates a symlink to `emails/static` and sets `EMAIL_PREVIEW=1` so templates use `/static/...` instead of `EMAIL_ASSETS_BASE_URL`. SMTP render still sets `EMAIL_ASSETS_BASE_URL` to the API public URL.

### 6. Docs

Update `docs/DESIGN_DOC.md` if templates, props, or asset paths change.
