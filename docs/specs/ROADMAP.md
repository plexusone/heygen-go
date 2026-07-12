# HeyGen Go SDK - Roadmap

## Version History

| Version | Status | Description |
|---------|--------|-------------|
| v0.1.0 | In Progress | Initial release with core SDK |
| v0.2.0 | Planned | Extended API coverage |
| v0.3.0 | Planned | Production hardening |

---

## v0.1.0 - Core SDK

**Target:** Initial release

### Features

- [x] Project setup and documentation
- [x] Core HTTP client with authentication
- [x] Retry logic with exponential backoff
- [x] Error types with status codes and request IDs
- [x] Avatar API v3 (list, get, list looks)
- [x] LiveAvatar session management (LITE mode)
- [x] Basic examples (list-avatars, liveavatar-session)
- [x] README documentation
- [x] MkDocs documentation site

### Remaining for v0.1.0

- [ ] **Unit tests for avatar package**
- [ ] **Unit tests for liveavatar package**
- [ ] **Unit tests for heygen core package**
- [ ] Push to GitHub
- [ ] CI/CD with GitHub Actions (build, lint, test)
- [ ] Voice API (list, get) - currently placeholder
- [ ] WebSocket streaming client for LiveAvatar audio

### API Coverage

| API | Endpoints | Status |
|-----|-----------|--------|
| Avatars | `/v3/avatars`, `/v3/avatars/{id}`, `/v3/avatars/{id}/looks` | ✅ Implemented |
| LiveAvatar | `/v1/sessions/token`, `/v1/sessions/start`, `/v1/sessions/stop` | ✅ Implemented |
| Voices | `/v2/voices` | Placeholder |
| Videos | `/v2/video/generate` | Placeholder |
| Streaming | Legacy `/v1/streaming.*` | Placeholder (deprecated) |

---

## v0.2.0 - Extended API

**Target:** After v0.1.0 stabilizes

### Features

- [ ] Video generation API
- [ ] Video status polling
- [ ] Voice listing API
- [ ] Template API
- [ ] Webhook handling
- [ ] Batch video generation
- [ ] Custom avatar support

### API Coverage

| API | Endpoints | Status |
|-----|-----------|--------|
| Videos | `video.generate`, `video_status.get`, `video.list`, `video.delete` | Planned |
| Voices | `voice.list` | Planned |
| Templates | `template.list`, `template.get` | Planned |
| Webhooks | Callback handling | Planned |

---

## v0.3.0 - Production Hardening

**Target:** Production readiness

### Features

- [ ] Comprehensive error handling
- [ ] Rate limit handling with backoff
- [ ] Connection pooling
- [ ] Metrics and observability
- [ ] Context propagation
- [ ] Graceful shutdown

### Quality

- [ ] 80%+ test coverage
- [ ] Integration test suite
- [ ] Performance benchmarks
- [ ] Security audit

---

## omni-livekit Integration

### Panel Discussion Agent

**Location:** `cmd/livekit-agent-panel/` in omni-livekit

| Task | Description | Status |
|------|-------------|--------|
| Coordinator | Turn-taking logic for multiple panelists | Planned |
| Panelist | Individual agent with personality/voice | Planned |
| Transcript | Shared conversation context | Planned |
| LiveAvatar Integration | Avatar rendering for panelists | Planned |

### Provider Milestones

| Milestone | Description | Status |
|-----------|-------------|--------|
| LiveAvatar Provider | HeyGen LiveAvatar in LiveKit rooms | Planned |
| Audio Routing | Agent TTS → LiveAvatar → LiveKit | Planned |
| Event Sync | Coordinate speaking events | Planned |
| Fallback | Static image when LiveAvatar unavailable | Planned |
| Multi-Avatar | Multiple avatars in panel discussions | Planned |

---

## Future Considerations

### Potential v0.4.0+ Features

| Feature | Description | Priority |
|---------|-------------|----------|
| Persona API | Create/manage personas | Medium |
| Knowledge API | Upload knowledge bases | Low |
| Photo Avatar | Generate avatars from photos | Low |
| Video Translation | Translate existing videos | Low |
| Interactive Avatar | Two-way conversations | Medium |

### Abstraction Layer

Future `omniavatar` package providing unified interface:

```go
type AvatarSession interface {
    Start(ctx context.Context) error
    PushAudio(frame AudioFrame) error
    Interrupt() error
    WaitForReady(ctx context.Context) error
    Close() error
}
```

Implementations:

- `omniavatar-heygen` - HeyGen LiveAvatar
- `omniavatar-tavus` - Tavus CVI
- `omniavatar-bithuman` - bitHuman
- `omniavatar-simli` - Simli
- `omniavatar-static` - Static image fallback

---

## Dependencies Tracking

| Dependency | Current | Notes |
|------------|---------|-------|
| github.com/livekit/protocol | v1.49.0 | LiveKit token generation |
| nhooyr.io/websocket | - | WebSocket client (future) |

---

## Release Checklist

### Pre-Release (v0.1.0)

- [x] Core features implemented
- [ ] Unit tests passing
- [x] golangci-lint clean
- [x] Documentation complete
- [ ] CHANGELOG.json created
- [x] Examples working

### Release

- [ ] Push to GitHub
- [ ] Set up GitHub Actions CI
- [ ] Tag version v0.1.0
- [ ] Generate CHANGELOG.md
- [ ] Create GitHub release
