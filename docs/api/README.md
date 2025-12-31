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
| 마켓 코드 조회 | GET | https://apidocs.bithumb.com/v2.1.5/reference/마켓코드-조회 |

### 시세 캔들 조회

| API명 | 메서드 | 링크 |
|-------|--------|------|
| 분(Minute) 캔들 | GET | https://apidocs.bithumb.com/v2.1.5/reference/분minute-캔들-1 |
| 일(Day) 캔들 | GET | https://apidocs.bithumb.com/v2.1.5/reference/일day-캔들 |
| 주(Week) 캔들 | GET | https://apidocs.bithumb.com/v2.1.5/reference/주week-캔들 |
| 월(Month) 캔들 | GET | https://apidocs.bithumb.com/v2.1.5/reference/월month-캔들 |

### 시세 체결 조회

| API명 | 메서드 | 링크 |
|-------|--------|------|
| 최근 체결 내역 | GET | https://apidocs.bithumb.com/v2.1.5/reference/최근-체결-내역 |

### 시세 현재가(Ticker) 조회

| API명 | 메서드 | 링크 |
|-------|--------|------|
| 현재가 정보 | GET | https://apidocs.bithumb.com/v2.1.5/reference/현재가-정보 |

### 시세 호가 정보(Orderbook) 조회

| API명 | 메서드 | 링크 |
|-------|--------|------|
| 호가 정보 조회 | GET | https://apidocs.bithumb.com/v2.1.5/reference/호가-정보-조회 |

### 서비스 정보

| API명 | 메서드 | 링크 |
|-------|--------|------|
| 경보제 | GET | https://apidocs.bithumb.com/v2.1.5/reference/경보제 |
| 공지사항 조회 | GET | https://apidocs.bithumb.com/v2.1.5/reference/공지사항-조회 |
| 입출금 수수료 조회 | GET | https://apidocs.bithumb.com/v2.1.5/reference/입출금-수수료-조회 |

---

## API 2.0 PRIVATE API

### 자산

| API명 | 메서드 | 링크 |
|-------|--------|------|
| 전체 계좌 조회 | GET | https://apidocs.bithumb.com/v2.1.5/reference/전체-계좌-조회 |

### 주문

| API명 | 메서드 | 링크 |
|-------|--------|------|
| 주문 가능 정보 | GET | https://apidocs.bithumb.com/v2.1.5/reference/주문-가능-정보 |
| 개별 주문 조회 | GET | https://apidocs.bithumb.com/v2.1.5/reference/개별-주문-조회 |
| 주문 리스트 조회 | GET | https://apidocs.bithumb.com/v2.1.5/reference/주문-리스트-조회 |
| 주문하기 [BETA] | POST | https://apidocs.bithumb.com/v2.1.5/reference/주문하기 |
| 주문 취소 접수 [BETA] | DELETE | https://apidocs.bithumb.com/v2.1.5/reference/주문-취소-접수 |

### 알고리즘 주문

#### TWAP

| API명 | 메서드 | 링크 |
|-------|--------|------|
| TWAP - 주문 내역 조회 | GET | https://apidocs.bithumb.com/v2.1.5/reference/twap-주문내역-조회 |
| TWAP - 주문 취소 | DELETE | https://apidocs.bithumb.com/v2.1.5/reference/twap-주문-취소 |
| TWAP - 주문하기 | POST | https://apidocs.bithumb.com/v2.1.5/reference/twap-주문-요청 |

### 출금

| API명 | 메서드 | 링크 |
|-------|--------|------|
| 코인 출금 리스트 조회 | GET | https://apidocs.bithumb.com/v2.1.5/reference/출금-리스트-조회 |
| 원화 출금 리스트 조회 | GET | https://apidocs.bithumb.com/v2.1.5/reference/원화-출금-리스트-조회 |
| 개별 출금 조회 | GET | https://apidocs.bithumb.com/v2.1.5/reference/개별-출금-조회 |
| 출금 가능 정보 | GET | https://apidocs.bithumb.com/v2.1.5/reference/출금-가능-정보 |
| 가상 자산 출금하기 | POST | https://apidocs.bithumb.com/v2.1.5/reference/디지털-자산-출금하기 |
| 원화 출금하기 | POST | https://apidocs.bithumb.com/v2.1.5/reference/원화-출금하기 |
| 출금 허용 주소 리스트 조회 | GET | https://apidocs.bithumb.com/v2.1.5/reference/출금-허용-주소-리스트-조회 |

### 입금

| API명 | 메서드 | 링크 |
|-------|--------|------|
| 코인 입금 리스트 조회 | GET | https://apidocs.bithumb.com/v2.1.5/reference/입금-리스트-조회 |
| 원화 입금 리스트 조회 | GET | https://apidocs.bithumb.com/v2.1.5/reference/원화-입금-리스트-조회 |
| 개별 입금 조회 | GET | https://apidocs.bithumb.com/v2.1.5/reference/개별-입금-조회 |
| 입금 주소 생성 요청 | POST | https://apidocs.bithumb.com/v2.1.5/reference/입금-주소-생성-요청 |
| 전체 입금 주소 조회 | GET | https://apidocs.bithumb.com/v2.1.5/reference/전체-입금-주소-조회 |
| 개별 입금 주소 조회 | GET | https://apidocs.bithumb.com/v2.1.5/reference/개별-입금-주소-조회 |
| 원화 입금하기 | POST | https://apidocs.bithumb.com/v2.1.5/reference/원화-입금하기 |

### 서비스 정보

| API명 | 메서드 | 링크 |
|-------|--------|------|
| 입출금 현황 | GET | https://apidocs.bithumb.com/v2.1.5/reference/입출금-현황 |
| API 키 리스트 조회 | GET | https://apidocs.bithumb.com/v2.1.5/reference/api-키-리스트-조회 |

### 코인대여(렌딩플러스)

| API명 | 메서드 | 링크 |
|-------|--------|------|
| 메이저 자산 조회(담보및거래가능) | GET | https://apidocs.bithumb.com/v2.1.5/reference/메이저-자산-조회담보및거래가능 |
| 상환레벨 조회 | GET | https://apidocs.bithumb.com/v2.1.5/reference/상환레벨-조회 |

---

## WEBSOCKET [Beta]

| 항목 | 링크 |
|------|------|
| 기본 정보 | https://apidocs.bithumb.com/v2.1.5/reference/기본-정보 |
| 요청 방법 및 포맷 | https://apidocs.bithumb.com/v2.1.5/reference/요청-포맷 |
| 테스트 및 요청 예제 | https://apidocs.bithumb.com/v2.1.5/reference/테스트-및-요청-예제 |

### 타입별 요청 및 응답

| 타입 | 링크 |
|------|------|
| 현재가 (Ticker) | https://apidocs.bithumb.com/v2.1.5/reference/현재가-ticker |
| 체결 (Trade) | https://apidocs.bithumb.com/v2.1.5/reference/체결-trade |
| 호가 (Orderbook) | https://apidocs.bithumb.com/v2.1.5/reference/호가-orderbook |
| 내 주문 및 체결 (MyOrder) | https://apidocs.bithumb.com/v2.1.5/reference/내-주문-및-체결-myorder |
| 내 자산 (MyAsset) | https://apidocs.bithumb.com/v2.1.5/reference/내-자산-myasset |

| 항목 | 링크 |
|------|------|
| 웹소켓 에러 | https://apidocs.bithumb.com/v2.1.5/reference/웹소켓-에러 |
| 연결 관리 | https://apidocs.bithumb.com/v2.1.5/reference/연결-관리 |

---

## 요약

| 구분 | API 수 |
|------|--------|
| PUBLIC API | 13개 |
| PRIVATE API | 28개 |
| WEBSOCKET 기본 | 8개 |
| **합계** | **49개** |

---

*마지막 업데이트: 2026-01-01*
*API 버전: v2.1.5*
*원본 문서: https://apidocs.bithumb.com/v2.1.5/reference*
