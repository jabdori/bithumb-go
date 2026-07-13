# 기여 가이드

Bithumb Go SDK에 관심을 가져주셔서 감사합니다! 기여를 환영하며, 이 가이드는 프로젝트에 기여하는 방법을 안내합니다.

## 목차

- [행동 강령](#행동-강령)
- [기여하는 방법](#기여하는-방법)
- [개발 환경 설정](#개발-환경-설정)
- [코드 스타일](#코드-스타일)
- [테스트](#테스트)
- [커밋 메시지](#커밋-메시지)
- [Pull Request](#pull-request)

## 행동 강령

- 존중과 배려를 가지고 상호작용해주세요
- 건설적인 피드백을 환영합니다
- 다양한 관점과 경험을 존중해주세요

## 기여하는 방법

### 버그 신고

버그를 발견하면 [GitHub Issues](https://github.com/hysuki/bithumb-go/issues)에 신고해주세요:

1. 기존 이슈를 검색하여 중복인지 확인
2. 제목에 `[Bug]` 접두사 추가
3. 다음 정보 포함:
   - SDK 버전
   - Go 버전 (`go version`)
   - OS 버전
   - 재현 가능한 최소 예제 코드
   - 예상 동작 vs 실제 동작
   - 에러 메시지 및 스택 트레이스

### 기능 요청

새로운 기능을 요청하려면:

1. 기존 이슈를 검색하여 중복인지 확인
2. 제목에 `[Feature]` 접두사 추가
3. 기능에 대한 상세한 설명
4. 사용 사례와 이유
5. 가능하다면 구현 제안

## 개발 환경 설정

### Fork 및 Clone

```bash
# 리포지토리 포크
# https://github.com/hysuki/bithumb-go/fork

# 로컬로 클론
git clone https://github.com/YOUR_USERNAME/bithumb-go.git
cd bithumb-go

# 업스트림 원격 저장소 추가
git remote add upstream https://github.com/hysuki/bithumb-go.git
```

### 의존성 설치

```bash
# 의존성 다운로드
go mod download

# 의존성 정리
go mod tidy
```

### 개발 브랜치 생성

```bash
# 메인 브랜치 업데이트
git checkout main
git pull upstream main

# 기능 브랜치 생성
git checkout -b feature/your-feature-name
# 또는 버그 수정
git checkout -b fix/your-bug-fix
```

## 코드 스타일

### Go 표준 준수

- [Effective Go](https://golang.org/doc/effective_go) 따르기
- `gofmt`로 코드 포맷팅
- `golint` 및 `go vet`으로 린팅

```bash
# 코드 포맷팅
gofmt -w .

# 린팅
golint ./...

# 정적 분석
go vet ./...
```

### 명명 규칙

```go
// 패키지 이름은 소문자, 한 단어
package public

// 상수는 UpperCamelCase
const (
    DefaultTimeout = 30 * time.Second
    MaxRetries     = 3
)

// 내보내지는 함수/타입은 UpperCamelCase
type Client struct {}
func NewClient() *Client {}

// 내보내지지 않은 것은 lowerCamelCase
type internalConfig struct {}
func handleError() error {}
```

### 주석

```go
// Package public provides a client for Bithumb Public API.
package public

// Client provides access to Bithumb Public API.
// Client values are safe for concurrent use.
type Client struct {
    // base is the underlying base client.
    base base.Client
}

// GetTicker retrieves current ticker information for the given markets.
// It returns an error if the request fails or if the response is invalid.
func (c *Client) GetTicker(req *GetTickerRequest) ([]*Ticker, error) {
    // Implementation...
}
```

### 에러 처리

```go
// 에러 래핑 (Go 1.13+)
if err != c.doRequest(ctx, req); err != nil {
    return fmt.Errorf("failed to get ticker: %w", err)
}

// 커스텀 에러 타입
var (
    ErrInvalidMarket = errors.New("invalid market")
    ErrEmptyResponse = errors.New("empty response")
)

// 에러 감지
if errors.Is(err, ErrInvalidMarket) {
    // Handle invalid market
}
```

## 테스트

### 단위 테스트 작성

```go
func TestClient_GetTicker(t *testing.T) {
    tests := []struct {
        name    string
        req     *GetTickerRequest
        want    []*Ticker
        wantErr bool
    }{
        {
            name: "valid request",
            req:  &GetTickerRequest{Markets: []string{"KRW-BTC"}},
            wantErr: false,
        },
        {
            name:    "empty markets",
            req:     &GetTickerRequest{Markets: []string{}},
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            c := NewTestClient()
            got, err := c.GetTicker(tt.req)
            if (err != nil) != tt.wantErr {
                t.Errorf("GetTicker() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            // Additional assertions...
        })
    }
}
```

### 테스트 실행

```bash
# 외부 API에 의존하지 않는 기본 테스트 실행
go test ./...

# 실제 빗썸 API를 호출하는 통합 테스트 실행
go test -tags=integration ./public

# 상세 모드
go test -v ./...

# 커버리지 확인
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# 특정 패키지 테스트
go test ./client
go test ./public -run TestGetMarketAll
go test -tags=integration ./public -run TestGetTicker

# 벤치마크
go test -bench=. -benchmem
```

### 테스트 커버리지

- 새로운 코드는 테스트 커버리지 80% 이상 권장
- 핵심 기능은 100% 커버리지 권장
- 테스트는 명확하고 유지보수 가능해야 함

### 모의(Mock) 사용

```go
// 내부 패키지에 모의 인터페이스
type mockHTTPClient struct {
    DoFunc func(*http.Request) (*http.Response, error)
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
    return m.DoFunc(req)
}
```

## 커밋 메시지

### 커밋 메시지 규칙

[Conventional Commits](https://www.conventionalcommits.org/) 따르기:

```
<type>(<scope>): <subject>

<body>

<footer>
```

### 타입

- `feat`: 새로운 기능
- `fix`: 버그 수정
- `docs`: 문서만 변경
- `style`: 코드 포맷팅, 세미콜론 누락 등 (코드 동작에 영향 없음)
- `refactor`: 코드 리팩토링 (새 기능이나 버그 수정 아님)
- `perf`: 성능 향상
- `test`: 테스트 추가/수정
- `chore`: 빌드 프로세스, 도구 설정 등

### 예시

```
feat(public): add GetCandlestick method

Implement candlestick data retrieval with support for
multiple time intervals (1, 3, 5, 10, 15, 30, 60, 240 minutes).

Closes #123
```

```
fix(auth): handle JWT expiration correctly

Previously, expired tokens were not refreshed properly,
causing authentication failures. This fix ensures tokens
are refreshed before expiration.

Fixes #456
```

```
docs: update API documentation with new endpoints

Add documentation for the newly added order management
endpoints including PlaceOrder and CancelOrder.
```

## Pull Request

### PR 전 체크리스트

- [ ] 코드가 `gofmt`으로 포맷팅됨
- [ ] `go vet ./...` 실행 후 에러 없음
- [ ] `go test ./...` 실행 후 모든 테스트 통과
- [ ] 새로운 기능에 대한 테스트 추가
- [ ] 문서 업데이트 (필요한 경우)
- [ ] 커밋 메시지가 명확함

### PR 제목

```
[type] Short description of changes
```

예시:
```
[Feature] Add WebSocket subscription support
[Fix] Handle API rate limiting correctly
[Docs] Update README with new examples
```

### PR 설명

```markdown
## 변경 사항
- 변경 사항 요약

## 관련 이슈
Closes #123

## 테스트 방법
```bash
# 테스트 단계
go test ./...
```

## 스크린샷/예시
(해당하는 경우)
```

### 코드 리뷰

- 모든 PR은 최소 1명의 승인자가 필요합니다
- CI/CD 파이프라인이 통과해야 합니다
- 리뷰어의 피드백을 적용하거나 의견을 명확히 하세요

### 변경 사항 요청

변경 사항이 요청되면:

1. 피드백을 확인하고 이해
2. 필요한 변경 사항 구현
3. 테스트 통과 확인
4. PR 상태를 "Ready for review"로 변경

## 릴리스 프로세스

1. 메인 브랜치로 머지
2. 버전 태그 생성 (예: `v1.2.0`)
3. CHANGELOG.md 업데이트
4. GitHub Release 생성

## 질문?

질문이 있으시면 [GitHub Issues](https://github.com/hysuki/bithumb-go/issues)에 등록해주세요.

---

다시 한번 기여에 감사드립니다! 🎉
