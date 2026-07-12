# HeyGen Go SDK - Product Requirements Document

## Overview

heygen-go is a Go SDK for the HeyGen API, providing programmatic access to HeyGen's video generation and LiveAvatar real-time streaming capabilities.

## Problem Statement

HeyGen provides powerful AI video generation and real-time avatar streaming, but lacks an official Go SDK. The PlexusOne ecosystem requires Go-native integration for:

- **OmniAgent**: Real-time voice AI with avatar rendering
- **omni-livekit**: LiveKit voice agents with visual representation
- **Video automation**: Programmatic video generation workflows

## Target Users

1. **PlexusOne internal** - Primary consumer via omni-livekit avatar providers
2. **Go developers** - Building applications with HeyGen's API
3. **DevOps/Automation** - CI/CD pipelines generating videos

## Goals

### Primary Goals

1. **REST API Coverage** - Full coverage of HeyGen's REST API
2. **LiveAvatar Support** - Real-time streaming with WebSocket protocol
3. **Provider Integration** - Enable omni-livekit HeyGen avatar provider
4. **Type Safety** - Generated types from OpenAPI spec

### Non-Goals (V1)

- Browser/frontend SDK (covered by @heygen/liveavatar-web-sdk)
- Video editing capabilities
- Avatar training/customization APIs

## Features

### V1.0 - Core SDK

| Feature | Priority | Description |
|---------|----------|-------------|
| Authentication | P0 | API key authentication |
| Avatars API | P0 | List, get avatar details |
| Voices API | P0 | List, get voice details |
| Video Generation | P1 | Create videos from scripts |
| LiveAvatar Sessions | P0 | Create, manage real-time sessions |
| WebSocket Streaming | P0 | Audio streaming to LiveAvatar |
| Session Events | P0 | Handle avatar events (ready, speaking, etc.) |

### V1.1 - Extended Features

| Feature | Priority | Description |
|---------|----------|-------------|
| Templates API | P2 | Video template management |
| Webhook Support | P2 | Video generation callbacks |
| Batch Operations | P2 | Bulk video generation |

## Success Metrics

1. **API Coverage** - 100% of documented REST endpoints
2. **LiveAvatar Parity** - Feature parity with Python livekit-plugins-liveavatar
3. **Integration** - Working omni-livekit HeyGen provider
4. **Documentation** - Complete usage examples and API reference

## Dependencies

| Dependency | Purpose |
|------------|---------|
| HeyGen API | Backend service |
| ogen | OpenAPI code generation |
| nhooyr.io/websocket | WebSocket client |
| omni-livekit | Avatar provider integration |

## Timeline

| Phase | Milestone | Target |
|-------|-----------|--------|
| 1 | OpenAPI generation + core client | Week 1 |
| 2 | LiveAvatar WebSocket integration | Week 2 |
| 3 | omni-livekit provider | Week 3 |
| 4 | Documentation + release | Week 4 |

## References

- [HeyGen API Reference](https://docs.heygen.com/reference)
- [HeyGen CLI (Go)](https://github.com/heygen-com/heygen-cli)
- [LiveAvatar Python Starter](https://github.com/heygen-com/liveavatar-starter-livekit-agent-python)
- [LiveKit LiveAvatar Plugin](https://docs.livekit.io/agents/models/avatar/plugins/liveavatar/)
