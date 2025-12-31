# Bithumb Go SDK Examples

이 디렉토리는 Bithumb Go SDK의 사용 예제들을 포함하고 있습니다.

## 예제 목록

### 1. Public API (`public_api/`)

API 키 없이 시장 데이터를 조회하는 예제입니다.

```bash
cd examples/public_api
go run main.go
```

**기능:**
- 현재가 (Ticker) 조회
- 호가 (OrderBook) 조회
- 체결 내역 조회
- 캔들 데이터 조회

### 2. Private API (`private_api/`)

API 키가 필요한 개인 계정 기능을 사용하는 예제입니다.

**환경변수 설정:**
```bash
export BITHUMB_ACCESS_KEY="your-access-key"
export BITHUMB_SECRET_KEY="your-secret-key"
```

**실행:**
```bash
cd examples/private_api
go run main.go
```

**기능:**
- 계정 정보 조회
- 특정 코인 잔고 조회
- 주문 생성 (주석 처리됨)

### 3. WebSocket (`websocket/`)

실시간 시장 데이터를 수신하는 예제입니다.

```bash
cd examples/websocket
go run main.go
```

**기능:**
- WebSocket 연결
- Ticker 실시간 구독
- OrderBook 실시간 구독
- 자동 재연결
- 그레이스풀 셧다운

### 4. Basic Usage (`basic_usage.go`)

기본 사용법을 보여주는 간단한 예제입니다.

```bash
go run examples/basic_usage.go
```

## 주의 사항

1. **Private API 예제**를 실행하려면 빗썸에서 API 키를 발급받아야 합니다.
2. **주문 기능**은 실제 거래가 발생하므로 주의해서 사용하세요.
3. 예제 코드의 주문 관련 부분은 주석 처리되어 있습니다.

## API 키 발급 방법

1. [빗썸](https://www.bithumb.com)에 로그인
2. 마이페이지 > API 관리
3. API 키 발급
4. 접근 권한 설정 (거래, 출금 등)
5. 생성된 API 키와 Secret Key를 환경변수에 설정

## 추가 정보

- [SDK 문서](../README.md)
- [빗썸 API 문서](https://apidocs.bithumb.com)
