# 빗썸 API 레퍼런스 목차

빗썸 Open API v2.1.5 문서 목차입니다.

- [API 2.0 PUBLIC API](#api-20-public-api)
- [API 2.0 PRIVATE API](#api-20-private-api)
- [WEBSOCKET Beta](#websocket-beta)

---

## API 2.0 PUBLIC API

### 시세 종목 조회

| API명 | 메서드 | 링크 |
|-------|--------|------|
| 마켓 코드 조회 | GET | [./public/마켓코드-조회.md](./public/마켓코드-조회.md) |

### 시세 캔들 조회

| API명 | 메서드 | 링크 |
|-------|--------|------|
| 분(Minute) 캔들 | GET | [./public/분minute-캔들.md](./public/분minute-캔들.md) |
| 일(Day) 캔들 | GET | [./public/일day-캔들.md](./public/일day-캔들.md) |
| 주(Week) 캔들 | GET | [./public/주week-캔들.md](./public/주week-캔들.md) |
| 월(Month) 캔들 | GET | [./public/월month-캔들.md](./public/월month-캔들.md) |

### 시세 체결 조회

| API명 | 메서드 | 링크 |
|-------|--------|------|
| 최근 체결 내역 | GET | [./public/최근-체결-내역.md](./public/최근-체결-내역.md) |

### 시세 현재가(Ticker) 조회

| API명 | 메서드 | 링크 |
|-------|--------|------|
| 현재가 정보 | GET | [./public/현재가-정보.md](./public/현재가-정보.md) |

### 시세 호가 정보(Orderbook) 조회

| API명 | 메서드 | 링크 |
|-------|--------|------|
| 호가 정보 조회 | GET | [./public/호가-정보-조회.md](./public/호가-정보-조회.md) |

### 서비스 정보

| API명 | 메서드 | 링크 |
|-------|--------|------|
| 경보제 | GET | [./public/경보제.md](./public/경보제.md) |
| 공지사항 조회 | GET | [./public/공지사항-조회.md](./public/공지사항-조회.md) |
| 입출금 수수료 조회 | GET | [./public/입출금-수수료-조회.md](./public/입출금-수수료-조회.md) |

---

## API 2.0 PRIVATE API

### 자산

| API명 | 메서드 | 링크 |
|-------|--------|------|
| 전체 계좌 조회 | GET | [./private/전체-계좌-조회.md](./private/전체-계좌-조회.md) |

### 주문

| API명 | 메서드 | 링크 |
|-------|--------|------|
| 주문 가능 정보 | GET | [./private/주문-가능-정보.md](./private/주문-가능-정보.md) |
| 개별 주문 조회 | GET | [./private/개별-주문-조회.md](./private/개별-주문-조회.md) |
| 주문 리스트 조회 | GET | [./private/주문-리스트-조회.md](./private/주문-리스트-조회.md) |
| 주문하기 [BETA] | POST | [./private/주문하기.md](./private/주문하기.md) |
| 주문 취소 접수 [BETA] | DELETE | [./private/주문-취소-접수.md](./private/주문-취소-접수.md) |

### 알고리즘 주문

#### TWAP

| API명 | 메서드 | 링크 |
|-------|--------|------|
| TWAP - 주문 내역 조회 | GET | [./private/TWAP-주문-내역-조회.md](./private/TWAP-주문-내역-조회.md) |
| TWAP - 주문 취소 | DELETE | [./private/TWAP-주문-취소.md](./private/TWAP-주문-취소.md) |
| TWAP - 주문하기 | POST | [./private/TWAP-주문하기.md](./private/TWAP-주문하기.md) |

### 출금

| API명 | 메서드 | 링크 |
|-------|--------|------|
| 코인 출금 리스트 조회 | GET | [./private/코인-출금-리스트-조회.md](./private/코인-출금-리스트-조회.md) |
| 원화 출금 리스트 조회 | GET | [./private/원화-출금-리스트-조회.md](./private/원화-출금-리스트-조회.md) |
| 개별 출금 조회 | GET | [./private/개별-출금-조회.md](./private/개별-출금-조회.md) |
| 출금 가능 정보 | GET | [./private/출금-가능-정보.md](./private/출금-가능-정보.md) |
| 가상 자산 출금하기 | POST | [./private/가상-자산-출금하기.md](./private/가상-자산-출금하기.md) |
| 원화 출금하기 | POST | [./private/원화-출금하기.md](./private/원화-출금하기.md) |
| 출금 허용 주소 리스트 조회 | GET | [./private/출금-허용-주소-리스트-조회.md](./private/출금-허용-주소-리스트-조회.md) |

### 입금

| API명 | 메서드 | 링크 |
|-------|--------|------|
| 코인 입금 리스트 조회 | GET | [./private/코인-입금-리스트-조회.md](./private/코인-입금-리스트-조회.md) |
| 원화 입금 리스트 조회 | GET | [./private/원화-입금-리스트-조회.md](./private/원화-입금-리스트-조회.md) |
| 개별 입금 조회 | GET | [./private/개별-입금-조회.md](./private/개별-입금-조회.md) |
| 입금 주소 생성 요청 | POST | [./private/입금-주소-생성-요청.md](./private/입금-주소-생성-요청.md) |
| 전체 입금 주소 조회 | GET | [./private/전체-입금-주소-조회.md](./private/전체-입금-주소-조회.md) |
| 개별 입금 주소 조회 | GET | [./private/개별-입금-주소-조회.md](./private/개별-입금-주소-조회.md) |
| 원화 입금하기 | POST | [./private/원화-입금하기.md](./private/원화-입금하기.md) |

### 서비스 정보

| API명 | 메서드 | 링크 |
|-------|--------|------|
| 입출금 현황 | GET | [./private/입출금-현황.md](./private/입출금-현황.md) |
| API 키 리스트 조회 | GET | [./private/API-키-리스트-조회.md](./private/API-키-리스트-조회.md) |

### 코인대여(렌딩플러스)

| API명 | 메서드 | 링크 |
|-------|--------|------|
| 메이저 자산 조회(담보및거래가능) | GET | [./private/메이저-자산-조회(담보및거래가능).md](./private/메이저-자산-조회(담보및거래가능).md) |
| 상환레벨 조회 | GET | [./private/상환레벨-조회.md](./private/상환레벨-조회.md) |

---

## WEBSOCKET [Beta]

| 항목 | 링크 |
|------|------|
| 기본 정보 | [./websocket/기본-정보.md](./websocket/기본-정보.md) |
| 요청 방법 및 포맷 | [./websocket/요청-방법-및-포맷.md](./websocket/요청-방법-및-포맷.md) |
| 테스트 및 요청 예제 | [./websocket/테스트-및-요청-예제.md](./websocket/테스트-및-요청-예제.md) |

### 타입별 요청 및 응답

| 타입 | 링크 |
|------|------|
| 현재가 (Ticker) | [./websocket/현재가-ticker.md](./websocket/현재가-ticker.md) |
| 체결 (Trade) | [./websocket/체결-trade.md](./websocket/체결-trade.md) |
| 호가 (Orderbook) | [./websocket/호가-orderbook.md](./websocket/호가-orderbook.md) |
| 내 주문 및 체결 (MyOrder) | [./websocket/내-주문-및-체결-myorder.md](./websocket/내-주문-및-체결-myorder.md) |
| 내 자산 (MyAsset) | [./websocket/내-자산-myasset.md](./websocket/내-자산-myasset.md) |

| 항목 | 링크 |
|------|------|
| 웹소켓 에러 | [./websocket/웹소켓-에러.md](./websocket/웹소켓-에러.md) |
| 연결 관리 | [./websocket/연결-관리.md](./websocket/연결-관리.md) |

---

## 요약

| 구분 | API 수 |
|------|--------|
| PUBLIC API | 11개 |
| PRIVATE API | 27개 |
| WEBSOCKET | 10개 |
| **합계** | **48개** |

---

*마지막 업데이트: 2026-01-01*
*API 버전: v2.1.5*
*원본 문서: https://apidocs.bithumb.com/v2.1.5/reference*
