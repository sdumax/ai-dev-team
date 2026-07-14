# Playwright MCP Guide

## Visual Verification Workflow

When verifying frontend features with Playwright MCP:

### 1. Start Dev Server
```bash
npm run dev  # or appropriate command
```

### 2. Take Screenshots
- Capture at breakpoints: 320px, 768px, 1024px, 1280px
- Verify visual consistency

### 3. Test Interactions
- Form submissions
- Button clicks and navigation
- Keyboard shortcuts
- Touch interactions

### 4. Accessibility Audit
- Run Playwright accessibility checks
- Verify ARIA labels and roles
- Test keyboard navigation
- Check color contrast (4.5:1 minimum)

### 5. Cross-Browser Testing
- Test in Chrome, Firefox, Safari
- Verify responsive behavior
- Check touch interactions on mobile

## Agent Responsibilities

| Agent | Playwright Usage |
|-------|-----------------|
| UI/UX Designer | Visual verification, accessibility audit, design compliance |
| Developer | Functional verification during implementation |
| QA | Functional testing, interaction testing |
| Reviewer | Visual verification during code review |

## Reference

When implementing frontend features, follow this guide instead of duplicating Playwright instructions in each agent file.
