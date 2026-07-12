# HeyGen Go SDK - Roadmap

## Version History

| Version | Status | Description |
|---------|--------|-------------|
| v0.1.0 | Planned | Initial release with core SDK |
| v0.2.0 | Planned | Extended API coverage |
| v0.3.0 | Planned | Production hardening |

---

## v0.1.0 - Core SDK

**Target:** Initial release

### Features

- [x] Project setup and documentation
- [ ] OpenAPI code generation
- [ ] Core HTTP client with authentication
- [ ] Avatar API (list, get)
- [ ] Voice API (list, get)
- [ ] LiveAvatar session management
- [ ] WebSocket streaming client
- [ ] Audio frame handling
- [ ] Event callbacks
- [ ] Basic examples
- [ ] README documentation

### API Coverage

| API | Endpoints | Status |
|-----|-----------|--------|
| Avatars | `avatar.list`, `avatar.get` | Planned |
| Voices | `voice.list`, `voice.get` | Planned |
| LiveAvatar | `streaming.new`, `streaming.start`, `streaming.stop` | Planned |
| WebSocket | Audio streaming, events | Planned |

### Integration

- [ ] omni-livekit HeyGen provider (separate PR in omni-livekit)

---

## v0.2.0 - Extended API

**Target:** After v0.1.0 stabilizes

### Features

- [ ] Video generation API
- [ ] Video status polling
- [ ] Template API
- [ ] Webhook handling
- [ ] Batch video generation
- [ ] Custom avatar support

### API Coverage

| API | Endpoints | Status |
|-----|-----------|--------|
| Videos | `video.generate`, `video_status.get` | Planned |
| Templates | `template.list`, `template.get` | Planned |
| Webhooks | Callback handling | Planned |

---

## v0.3.0 - Production Hardening

**Target:** Production readiness

### Features

- [ ] Comprehensive error handling
- [ ] Retry with exponential backoff
- [ ] Rate limit handling
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

## Future Considerations

### Potential v0.4.0+ Features

| Feature | Description | Priority |
|---------|-------------|----------|
| Persona API | Create/manage personas | Medium |
| Knowledge API | Upload knowledge bases | Low |
| Photo Avatar | Generate avatars from photos | Low |
| Video Translation | Translate existing videos | Low |
| Interactive Avatar | Two-way conversations | Medium |

### omni-livekit Integration Milestones

| Milestone | Description |
|-----------|-------------|
| Provider MVP | Basic HeyGen avatar in LiveKit rooms |
| Audio Routing | Agent TTS → HeyGen → LiveKit |
| Event Sync | Coordinate speaking events |
| Fallback | Static image when HeyGen unavailable |
| Multi-Avatar | Multiple HeyGen avatars in panel discussions |

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

| Dependency | Current | Latest | Notes |
|------------|---------|--------|-------|
| ogen | - | latest | OpenAPI generation |
| nhooyr.io/websocket | - | v1.8.x | WebSocket client |
| omni-livekit | v0.3.0 | v0.3.0 | Avatar provider integration |

---

## Release Checklist

### Pre-Release

- [ ] All planned features implemented
- [ ] Tests passing
- [ ] golangci-lint clean
- [ ] Documentation complete
- [ ] CHANGELOG.json updated
- [ ] Examples working

### Release

- [ ] Tag version
- [ ] Generate CHANGELOG.md
- [ ] Create GitHub release
- [ ] Update omni-livekit dependency
