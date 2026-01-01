# TWAP - 주문 내역 조회

> **원본 URL**: https://apidocs.bithumb.com/v2.1.5/reference/twap-주문내역-조회

## 설명

TWAP 주문 목록을 조회합니다.

## Response

| 필드 | 설명 | 타입 |
|------|------|------|
| market | Market ID (ex.KRW-BTC) | String |
| uuids | TWAP 주문 ID 목록 | Array |
| state | TWAP 주문 상태
- progress : 진행중 (default)
- done : 주문 완료
- cancel : 취소 | String |
| next_key | 다음 페이지 조회를 위한 커서 값 | String |
| limit | Number limit (default: 100, limit: 100) | Number |
| order_by | Sorting method
- asc : 오름차순
- desc : 내림차순 (default) | String |

