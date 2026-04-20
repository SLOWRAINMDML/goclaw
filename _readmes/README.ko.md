<p align="center">
  <img src="../_statics/goclaw.png" alt="GoClaw" />
</p>

<h1 align="center">GoClaw</h1>

<p align="center"><strong>Enterprise AI Agent Platform</strong></p>

<p align="center">
Multi-agent AI gateway built in Go. 20+ LLM providers. 7 channels. Multi-tenant PostgreSQL.<br/>
Single binary. Production-tested. Agents that orchestrate for you.
</p>

<p align="center">
  <a href="https://docs.goclaw.sh">문서</a> •
  <a href="https://docs.goclaw.sh/#quick-start">빠른 시작</a> •
  <a href="https://x.com/nlb_io">Twitter / X</a>
</p>

<p align="center">
  <a href="https://go.dev/"><img src="https://img.shields.io/badge/Go_1.26-00ADD8?style=flat-square&logo=go&logoColor=white" alt="Go" /></a>
  <a href="https://www.postgresql.org/"><img src="https://img.shields.io/badge/PostgreSQL_18-316192?style=flat-square&logo=postgresql&logoColor=white" alt="PostgreSQL" /></a>
  <a href="https://www.docker.com/"><img src="https://img.shields.io/badge/Docker-2496ED?style=flat-square&logo=docker&logoColor=white" alt="Docker" /></a>
  <a href="https://developer.mozilla.org/en-US/docs/Web/API/WebSocket"><img src="https://img.shields.io/badge/WebSocket-010101?style=flat-square&logo=socket.io&logoColor=white" alt="WebSocket" /></a>
  <a href="https://opentelemetry.io/"><img src="https://img.shields.io/badge/OpenTelemetry-000000?style=flat-square&logo=opentelemetry&logoColor=white" alt="OpenTelemetry" /></a>
  <a href="https://www.anthropic.com/"><img src="https://img.shields.io/badge/Anthropic-191919?style=flat-square&logo=anthropic&logoColor=white" alt="Anthropic" /></a>
  <a href="https://openai.com/"><img src="https://img.shields.io/badge/OpenAI_Compatible-412991?style=flat-square&logo=openai&logoColor=white" alt="OpenAI" /></a>
  <img src="https://img.shields.io/badge/License-MIT-yellow?style=flat-square" alt="License: MIT" />
</p>

**GoClaw**는 LLM을 도구, 채널, 데이터에 연결하는 멀티 에이전트 AI 게이트웨이입니다 — 런타임 의존성 없이 단일 Go 바이너리로 배포됩니다. 20개 이상의 LLM 공급자에서 완전한 멀티 테넌트 격리와 함께 에이전트 팀과 에이전트 간 위임을 조율합니다.

향상된 보안, 멀티 테넌트 PostgreSQL, 프로덕션 수준의 관측 가능성을 갖춘 [OpenClaw](https://github.com/openclaw/openclaw)의 Go 포트입니다.

🌐 **Languages:**
[🇺🇸 English](../README.md) ·
[🇨🇳 简体中文](README.zh-CN.md) ·
[🇯🇵 日本語](README.ja.md) ·
[🇰🇷 한국어](README.ko.md) ·
[🇻🇳 Tiếng Việt](README.vi.md) ·
[🇪🇸 Español](README.es.md) ·
[🇧🇷 Português](README.pt.md) ·
[🇫🇷 Français](README.fr.md) ·
[🇩🇪 Deutsch](README.de.md) ·
[🇷🇺 Русский](README.ru.md)
## 차별점

- **에이전트 팀 & 오케스트레이션** — 공유 태스크 보드, 에이전트 간 위임(동기/비동기), 하이브리드 에이전트 디스커버리를 갖춘 팀
- **멀티 테넌트 PostgreSQL** — 사용자별 워크스페이스, 사용자별 컨텍스트 파일, 암호화된 API 키(AES-256-GCM), 격리된 세션
- **단일 바이너리** — 약 25MB 정적 Go 바이너리, Node.js 런타임 불필요, 1초 미만 시작, $5 VPS에서 실행 가능
- **프로덕션 보안** — 5계층 권한 시스템(게이트웨이 인증 → 전역 도구 정책 → 에이전트별 → 채널별 → 소유자 전용)과 속도 제한, 프롬프트 인젝션 감지, SSRF 보호, 셸 차단 패턴, AES-256-GCM 암호화
- **20개 이상의 LLM 공급자** — Anthropic(네이티브 HTTP+SSE, 프롬프트 캐싱), OpenAI, OpenRouter, Groq, DeepSeek, Gemini, Mistral, xAI, MiniMax, Cohere, Perplexity, DashScope, Bailian, Zai, Ollama, Ollama Cloud, Claude CLI, Codex, ACP, 그리고 모든 OpenAI 호환 엔드포인트
- **7개 메시징 채널** — Telegram, Discord, Slack, Zalo OA, Zalo Personal, Feishu/Lark, WhatsApp
- **Extended Thinking** — 공급자별 사고 모드(Anthropic 예산 토큰, OpenAI 추론 노력, DashScope 사고 예산)와 스트리밍 지원
- **Heartbeat** — HEARTBEAT.md 체크리스트를 통한 주기적 에이전트 체크인(정상 시 억제, 활성 시간, 재시도 로직, 채널 전달)
- **스케줄링 & Cron** — 레인 기반 동시성으로 자동화된 에이전트 작업을 위한 `at`, `every`, cron 표현식
- **관측 가능성** — 스팬과 프롬프트 캐시 메트릭이 포함된 내장 LLM 호출 추적, 선택적 OpenTelemetry OTLP 내보내기

## Claw 에코시스템

|                 | OpenClaw        | ZeroClaw | PicoClaw | **GoClaw**                              |
| --------------- | --------------- | -------- | -------- | --------------------------------------- |
| 언어            | TypeScript      | Rust     | Go       | **Go**                                  |
| 바이너리 크기   | 28 MB + Node.js | 3.4 MB   | ~8 MB    | **~25 MB** (기본) / **~36 MB** (+ OTel) |
| Docker 이미지   | —               | —        | —        | **~50 MB** (Alpine)                     |
| RAM (유휴)      | > 1 GB          | < 5 MB   | < 10 MB  | **~35 MB**                              |
| 시작 시간       | > 5 s           | < 10 ms  | < 1 s    | **< 1 s**                               |
| 대상 하드웨어   | $599+ Mac Mini  | $10 엣지 | $10 엣지 | **$5 VPS+**                             |

| 기능                       | OpenClaw                             | ZeroClaw                                     | PicoClaw                              | **GoClaw**                     |
| -------------------------- | ------------------------------------ | -------------------------------------------- | ------------------------------------- | ------------------------------ |
| 멀티 테넌트 (PostgreSQL)   | —                                    | —                                            | —                                     | ✅                             |
| MCP 통합                   | — (ACP 사용)                         | —                                            | —                                     | ✅ (stdio/SSE/streamable-http) |
| 에이전트 팀                | —                                    | —                                            | —                                     | ✅ 태스크 보드 + 메일박스      |
| 보안 강화                  | ✅ (SSRF, 경로 순회, 인젝션)         | ✅ (샌드박스, 속도 제한, 인젝션, 페어링)     | 기본 (워크스페이스 제한, exec 차단)   | ✅ 5계층 방어                  |
| OTel 관측 가능성           | ✅ (옵트인 확장)                     | ✅ (Prometheus + OTLP)                       | —                                     | ✅ OTLP (옵트인 빌드 태그)     |
| 프롬프트 캐싱              | —                                    | —                                            | —                                     | ✅ Anthropic + OpenAI 호환     |
| 지식 그래프                | —                                    | —                                            | —                                     | ✅ LLM 추출 + 순회             |
| 스킬 시스템                | ✅ 임베딩/시맨틱                     | ✅ SKILL.md + TOML                           | ✅ 기본                               | ✅ BM25 + pgvector 하이브리드  |
| 레인 기반 스케줄러         | ✅                                   | 제한된 동시성                                | —                                     | ✅ (main/subagent/team/cron)   |
| 메시징 채널                | 37+                                  | 15+                                          | 10+                                   | 7+                             |
| 동반 앱                    | macOS, iOS, Android                  | Python SDK                                   | —                                     | 웹 대시보드                    |
| 라이브 캔버스 / 음성       | ✅ (A2UI + TTS/STT)                  | —                                            | 음성 전사                             | TTS (4개 공급자)               |
| LLM 공급자                 | 10+                                  | 8 네이티브 + 29 호환                         | 13+                                   | **20+**                        |
| 사용자별 워크스페이스      | ✅ (파일 기반)                       | —                                            | —                                     | ✅ (PostgreSQL)                |
| 암호화된 시크릿            | — (환경 변수만)                      | ✅ ChaCha20-Poly1305                         | — (평문 JSON)                         | ✅ DB의 AES-256-GCM            |

## 아키텍처

<p align="center">
  <img src="../_statics/architecture.jpg" alt="GoClaw Architecture" width="800" />
</p>

## 빠른 시작

**사전 요구사항:** Go 1.26+, pgvector가 포함된 PostgreSQL 18, Docker (선택 사항)

### 소스에서 빌드

```bash
git clone https://github.com/nextlevelbuilder/goclaw.git && cd goclaw
make build
./goclaw onboard        # Interactive setup wizard
source .env.local && ./goclaw
```

### 계층형 하위 워크스페이스 (WIP)

> **상태:** 실험적 스캐폴딩이 **2026-04-20 14:21:29 UTC** 기준 `dev` 브랜치에 반영됨  
> **최신 구현 커밋:** `fd1debb4` — `feat: add hierarchical workspace and skill runtime scaffolding`

현재 `dev` 브랜치에는 **마스터/관리자 워크스페이스 → 하위 워크스페이스** 운영 모델을 위한 첫 번째 기반 작업이 들어가 있습니다.
이 구현은 즉시 스키마 마이그레이션을 강제하지 않도록 **child tenant + tenant settings** 방식으로 시작했습니다.

#### 현재 가능한 것

- 마스터/관리자 tenant가 API를 통해 **하위 워크스페이스**를 생성할 수 있습니다.
- 하위 워크스페이스는 **스토리지, 워크스페이스 경로, 스킬 디렉토리**를 독립적으로 유지합니다.
- `inherit_parent_access`를 통해 부모 권한 상속을 선택적으로 허용할 수 있습니다.
- 각 하위 워크스페이스는 다음 메타데이터를 가질 수 있습니다.
  - git 연동 메타데이터
  - 채널 바인딩 메타데이터
  - 출력 바인딩 메타데이터
  - 코딩 전용 provider / model 선호값
  - reasoning 출력 모드 (`full`, `summary`, `none`)
- 공유/public 스킬을 **하위 워크스페이스로 포크**할 수 있습니다.
- 포크되거나 업로드된 스킬은 `.runtime/` 아래에 **가벼운 Python 런타임 래퍼**를 준비할 수 있습니다.
- 스킬 오작동 / 예시 / 이슈 피드백을 스킬 폴더에 구조화된 기록으로 남길 수 있습니다.

#### 이번 업데이트에서 추가된 엔드포인트

| 엔드포인트 | 목적 | 비고 |
|---|---|---|
| `GET /v1/tenants/{id}/children` | 부모 tenant의 하위 워크스페이스 목록 조회 | owner/admin 경로 |
| `POST /v1/tenants/{id}/children` | 하위 워크스페이스 생성 | child tenant settings로 저장 |
| `POST /v1/skills/{id}/fork` | 공유/public 스킬을 현재 tenant 공간으로 포크 | `.runtime`, `feedback`, `.env`, 명백한 credential 파일 제외 |
| `POST /v1/skills/{id}/prepare-runtime` | 스킬용 런타임 산출물 준비 | 현재 MVP는 Python 중심 |
| `GET /v1/skills/{id}/feedback` | 기록된 이슈/예시 피드백 조회 | `feedback/examples.jsonl` 읽기 |
| `POST /v1/skills/{id}/feedback` | 이슈/예시 피드백 추가 | 구조화된 JSONL 레코드 |

#### 현재 인식하는 child workspace 설정

| 설정 | 의미 |
|---|---|
| `parent_tenant_id` | 부모/마스터 워크스페이스 식별자 |
| `inherit_parent_access` | 부모 admin/owner가 child workspace를 접근/관리할 수 있도록 허용 |
| `workspace_mode` | 격리 전략 의도값. 현재 기본값은 `isolated` |
| `git_links` | 향후 sync/mirroring에 사용할 외부 git 저장소 메타데이터 |
| `channel_bindings` | 어떤 채널/chat/topic이 어떤 agent에 묶일지 설명하는 메타데이터 |
| `output_bindings` | webhook/API 채널 같은 outbound 대상 메타데이터 |
| `coding_provider` | 향후 코딩 작업 전용 라우팅에 사용할 선호 provider |
| `coding_model` | 코딩 작업용 선호 모델 |
| `reasoning_output` | 채널에 보이는 reasoning 정책: `full`, `summary`, `none` |

#### Reasoning 출력 모드

| 모드 | 동작 |
|---|---|
| `full` | 지원되는 스트리밍 채널에 reasoning preview를 그대로 표시 |
| `summary` | reasoning을 짧은 결정적(summary) 형태로 표시 |
| `none` | 채널에는 reasoning을 숨기고 최종 답변 흐름만 유지 |

#### 이번 업데이트에서 추가된 스킬 워크스페이스 동작

| 기능 | 현재 동작 |
|---|---|
| 스킬 포크 | 스킬을 현재 tenant 스킬 저장소로 복사하여 로컬 수정 가능 |
| 민감 파일 회피 | 포크 시 `.runtime`, `feedback`, `.env`, `secrets/`, 명백한 credential 계열 파일 제외 |
| 런타임 준비 | Python 엔트리포인트 스킬에 대해 `.runtime/run.py` + `.runtime/manifest.json` 생성 |
| 피드백 기록 | 구조화된 예시/이슈를 `feedback/examples.jsonl`에 append |

#### 현재 한계

이 구현은 **완성본이 아니라 방향을 잡기 위한 스캐폴딩**입니다.

- `channel_bindings`, `output_bindings`는 현재 **메타데이터로 저장만** 되고, 완전한 DB 기반 실시간 라우팅까지는 아직 연결되지 않았습니다.
- `git_links`는 현재 **설명용 메타데이터**이며, 실제 저장소 sync/mirroring 실행 로직은 아직 필요합니다.
- 부모/자식 권한 상속은 HTTP tenant 해석 쪽에 일부 연결됐지만, 전체 엔드투엔드 검증은 더 필요합니다.
- 런타임 준비는 현재 **Python MVP만** 대상으로 합니다.
- 전체 compile/build 검증은 Go 실행 가능한 환경에서 추가 확인이 필요합니다.

#### Child workspace 생성 예시 payload

```json
{
  "name": "client-a-workspace",
  "slug": "client-a",
  "inherit_parent_access": true,
  "workspace_mode": "isolated",
  "reasoning_output": "summary",
  "coding_provider": "codex",
  "coding_model": "gpt-5-codex",
  "git_links": [
    {
      "name": "client-a-repo",
      "repo_url": "https://github.com/example/client-a",
      "branch": "main",
      "directory": "workspace",
      "mode": "mirror"
    }
  ],
  "channel_bindings": [
    {
      "channel_instance_id": "<channel-instance-uuid>",
      "match_chat_ids": ["1234567890"],
      "agent_id": "<agent-uuid>",
      "enabled": true
    }
  ],
  "output_bindings": [
    {
      "type": "webhook",
      "target": "https://example.com/hooks/reply",
      "channel": "discord"
    }
  ]
}
```

#### 스킬 포크 플로우 예시

```bash
# 1) 공유/public 스킬을 현재 tenant 워크스페이스로 포크
curl -X POST http://localhost:18790/v1/skills/<skill-id>/fork \
  -H 'Authorization: Bearer <token>' \
  -H 'Content-Type: application/json' \
  -d '{
    "name": "client-a-custom-skill",
    "slug": "client-a-custom-skill",
    "visibility": "internal"
  }'

# 2) Python 엔트리포인트가 있다면 runtime 산출물 준비
curl -X POST http://localhost:18790/v1/skills/<forked-skill-id>/prepare-runtime \
  -H 'Authorization: Bearer <token>'

# 3) 스킬에 이슈/예시 피드백 기록
curl -X POST http://localhost:18790/v1/skills/<forked-skill-id>/feedback \
  -H 'Authorization: Bearer <token>' \
  -H 'Content-Type: application/json' \
  -d '{
    "type": "issue",
    "behavior_ref": "tool invocation",
    "code_path": "scripts/run.py",
    "prompt": "저장소를 요약해줘",
    "expected_behavior": "저장소의 핵심을 간결하게 요약",
    "actual_behavior": "누락된 dependency 오류 발생",
    "note": "dependency 부재 시 fallback 처리 필요"
  }'
```

### Docker 사용

```bash
# Generate .env with auto-generated secrets
chmod +x prepare-env.sh && ./prepare-env.sh

# Add at least one GOCLAW_*_API_KEY to .env, then:
make up

# Web Dashboard at http://localhost:18790
# Health check: curl http://localhost:18790/health
```

`GOCLAW_*_API_KEY` 환경 변수가 설정되면, 게이트웨이는 대화형 프롬프트 없이 자동으로 온보딩됩니다 — 공급자를 감지하고, 마이그레이션을 실행하며, 기본 데이터를 시드합니다.

> 빌드 변형(OTel, Tailscale, Redis), Docker 이미지 태그, compose 오버레이에 대해서는 [배포 가이드](https://docs.goclaw.sh/#deploy-docker-compose)를 참조하세요.

## 멀티 에이전트 오케스트레이션

GoClaw는 에이전트 팀과 에이전트 간 위임을 지원합니다 — 각 에이전트는 자체 ID, 도구, LLM 공급자, 컨텍스트 파일로 실행됩니다.

### 에이전트 위임

<p align="center">
  <img src="../_statics/agent-delegation.jpg" alt="Agent Delegation" width="700" />
</p>

| 모드 | 동작 방식 | 적합한 경우 |
|------|-------------|----------|
| **동기(Sync)** | 에이전트 A가 에이전트 B에게 요청하고 답변을 **기다림** | 빠른 조회, 사실 확인 |
| **비동기(Async)** | 에이전트 A가 에이전트 B에게 요청하고 **계속 진행**. B가 나중에 알림 | 장시간 작업, 보고서, 심층 분석 |

에이전트는 방향 제어(`outbound`, `inbound`, `bidirectional`)와 링크별 및 에이전트별 동시성 제한이 있는 명시적인 **권한 링크**를 통해 통신합니다.

### 에이전트 팀

<p align="center">
  <img src="../_statics/agent-teams.jpg" alt="Agent Teams Workflow" width="800" />
</p>

- **공유 태스크 보드** — `blocked_by` 의존성이 있는 태스크 생성, 클레임, 완료, 검색
- **팀 메일박스** — 직접 P2P 메시지 및 브로드캐스트
- **도구**: 태스크 관리를 위한 `team_tasks`, 메일박스를 위한 `team_message`

> 위임 세부 사항, 권한 링크, 동시성 제어에 대해서는 [에이전트 팀 문서](https://docs.goclaw.sh/#teams-what-are-teams)를 참조하세요.

## 내장 도구

| 도구               | 그룹          | 설명                                                         |
| ------------------ | ------------- | ------------------------------------------------------------ |
| `read_file`        | fs            | 파일 내용 읽기 (가상 FS 라우팅 포함)                         |
| `write_file`       | fs            | 파일 쓰기/생성                                               |
| `edit_file`        | fs            | 기존 파일에 대상 편집 적용                                   |
| `list_files`       | fs            | 디렉토리 내용 나열                                           |
| `search`           | fs            | 패턴으로 파일 내용 검색                                      |
| `glob`             | fs            | glob 패턴으로 파일 찾기                                      |
| `exec`             | runtime       | 셸 명령 실행 (승인 워크플로우 포함)                          |
| `web_search`       | web           | 웹 검색 (Brave, DuckDuckGo)                                  |
| `web_fetch`        | web           | 웹 콘텐츠 가져오기 및 파싱                                   |
| `memory_search`    | memory        | 장기 기억 검색 (FTS + 벡터)                                  |
| `memory_get`       | memory        | 기억 항목 검색                                               |
| `skill_search`     | —             | 스킬 검색 (BM25 + 임베딩 하이브리드)                         |
| `knowledge_graph_search` | memory  | 엔티티 검색 및 지식 그래프 관계 순회                         |
| `create_image`     | media         | 이미지 생성 (DashScope, MiniMax)                             |
| `create_audio`     | media         | 오디오 생성 (OpenAI, ElevenLabs, MiniMax, Suno)              |
| `create_video`     | media         | 비디오 생성 (MiniMax, Veo)                                   |
| `read_document`    | media         | 문서 읽기 (Gemini File API, 공급자 체인)                     |
| `read_image`       | media         | 이미지 분석                                                  |
| `read_audio`       | media         | 오디오 전사 및 분석                                          |
| `read_video`       | media         | 비디오 분석                                                  |
| `message`          | messaging     | 채널에 메시지 전송                                           |
| `tts`              | —             | Text-to-Speech 합성                                          |
| `spawn`            | —             | 서브에이전트 생성                                            |
| `subagents`        | sessions      | 실행 중인 서브에이전트 제어                                  |
| `team_tasks`       | teams         | 공유 태스크 보드 (나열, 생성, 클레임, 완료, 검색)            |
| `team_message`     | teams         | 팀 메일박스 (전송, 브로드캐스트, 읽기)                       |
| `sessions_list`    | sessions      | 활성 세션 나열                                               |
| `sessions_history` | sessions      | 세션 기록 보기                                               |
| `sessions_send`    | sessions      | 세션에 메시지 전송                                           |
| `sessions_spawn`   | sessions      | 새 세션 생성                                                 |
| `session_status`   | sessions      | 세션 상태 확인                                               |
| `cron`             | automation    | cron 작업 스케줄링 및 관리                                   |
| `gateway`          | automation    | 게이트웨이 관리                                              |
| `browser`          | ui            | 브라우저 자동화 (탐색, 클릭, 타이핑, 스크린샷)               |
| `announce_queue`   | automation    | 비동기 결과 알림 (비동기 위임용)                             |

## 문서

전체 문서는 **[docs.goclaw.sh](https://docs.goclaw.sh)**에서 — 또는 [`goclaw-docs/`](https://github.com/nextlevelbuilder/goclaw-docs)의 소스를 탐색하세요.

| 섹션 | 주제 |
|---------|--------|
| [시작하기](https://docs.goclaw.sh/#what-is-goclaw) | 설치, 빠른 시작, 구성, 웹 대시보드 둘러보기 |
| [핵심 개념](https://docs.goclaw.sh/#how-goclaw-works) | 에이전트 루프, 세션, 도구, 기억, 멀티 테넌시 |
| [에이전트](https://docs.goclaw.sh/#creating-agents) | 에이전트 생성, 컨텍스트 파일, 개성, 공유 & 접근 |
| [공급자](https://docs.goclaw.sh/#providers-overview) | Anthropic, OpenAI, OpenRouter, Gemini, DeepSeek, +15개 이상 |
| [채널](https://docs.goclaw.sh/#channels-overview) | Telegram, Discord, Slack, Feishu, Zalo, WhatsApp, WebSocket |
| [에이전트 팀](https://docs.goclaw.sh/#teams-what-are-teams) | 팀, 태스크 보드, 메시징, 위임 & 핸드오프 |
| [고급](https://docs.goclaw.sh/#custom-tools) | 커스텀 도구, MCP, 스킬, Cron, 샌드박스, 훅, RBAC |
| [배포](https://docs.goclaw.sh/#deploy-docker-compose) | Docker Compose, 데이터베이스, 보안, 관측 가능성, Tailscale |
| [레퍼런스](https://docs.goclaw.sh/#cli-commands) | CLI 명령, REST API, WebSocket 프로토콜, 환경 변수 |

## 테스트

```bash
go test ./...                                    # Unit tests
go test -v ./tests/integration/ -timeout 120s    # Integration tests (requires running gateway)
```

## 프로젝트 상태

프로덕션에서 테스트된 내용과 진행 중인 내용을 포함한 상세 기능 상태는 [CHANGELOG.md](CHANGELOG.md)를 참조하세요.

## 감사의 말

GoClaw는 원본 [OpenClaw](https://github.com/openclaw/openclaw) 프로젝트를 기반으로 만들어졌습니다. 이 Go 포트에 영감을 준 아키텍처와 비전에 감사드립니다.

## 라이선스

MIT
