# TWAP - 주문하기

> **원본 URL**: https://apidocs.bithumb.com/v2.1.5/reference/twap-주문-요청

## 설명

TWAP 주문을 요청합니다.

## Response

| 필드 | 설명 | 타입 |
|------|------|------|
| market * | Market ID (ex.KRW-BTC) | String |
| side * | Order Type
- bid : 매수
- ask : 매도 | String |
| volume | 주문량 (매도 시 필수) | NumberString |
| price | 주문 가격. (매수 시 필수) | NumberString |
| duration * | 주문 시간 (twap 주문이 진행되는 시간) - 초
- Min 300, Max 43200 | NumberString |
| frequency* | 주문 간격 - 초
- 5, 15, 20, 30, 60, 120 값 중에 하나 입력 가능 | NumberString |

