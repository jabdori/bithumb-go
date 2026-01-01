const { chromium } = require('playwright');
const fs = require('fs').promises;
const path = require('path');

const API_LIST = [
  { url: 'https://apidocs.bithumb.com/v2.1.5/reference/전체-계좌-조회', slug: '전체-계좌-조회', file: '전체-계좌-조회.md' },
  { url: 'https://apidocs.bithumb.com/v2.1.5/reference/주문-가능-정보', slug: '주문-가능-정보', file: '주문-가능-정보.md' },
  { url: 'https://apidocs.bithumb.com/v2.1.5/reference/개별-주문-조회', slug: '개별-주문-조회', file: '개별-주문-조회.md' },
  { url: 'https://apidocs.bithumb.com/v2.1.5/reference/주문-리스트-조회', slug: '주문-리스트-조회', file: '주문-리스트-조회.md' },
  { url: 'https://apidocs.bithumb.com/v2.1.5/reference/주문하기', slug: '주문하기', file: '주문하기.md' },
  { url: 'https://apidocs.bithumb.com/v2.1.5/reference/주문-취소-접수', slug: '주문-취소-접수', file: '주문-취소-접수.md' },
  { url: 'https://apidocs.bithumb.com/v2.1.5/reference/twap-주문내역-조회', slug: 'twap-주문내역-조회', file: 'TWAP-주문-내역-조회.md' },
  { url: 'https://apidocs.bithumb.com/v2.1.5/reference/twap-주문-취소', slug: 'twap-주문-취소', file: 'TWAP-주문-취소.md' },
  { url: 'https://apidocs.bithumb.com/v2.1.5/reference/twap-주문-요청', slug: 'twap-주문-요청', file: 'TWAP-주문하기.md' },
  { url: 'https://apidocs.bithumb.com/v2.1.5/reference/출금-리스트-조회', slug: '출금-리스트-조회', file: '코인-출금-리스트-조회.md' },
  { url: 'https://apidocs.bithumb.com/v2.1.5/reference/원화-출금-리스트-조회', slug: '원화-출금-리스트-조회', file: '원화-출금-리스트-조회.md' },
  { url: 'https://apidocs.bithumb.com/v2.1.5/reference/개별-출금-조회', slug: '개별-출금-조회', file: '개별-출금-조회.md' },
  { url: 'https://apidocs.bithumb.com/v2.1.5/reference/출금-가능-정보', slug: '출금-가능-정보', file: '출금-가능-정보.md' },
  { url: 'https://apidocs.bithumb.com/v2.1.5/reference/디지털-자산-출금하기', slug: '디지털-자산-출금하기', file: '가상-자산-출금하기.md' },
  { url: 'https://apidocs.bithumb.com/v2.1.5/reference/원화-출금하기', slug: '원화-출금하기', file: '원화-출금하기.md' },
  { url: 'https://apidocs.bithumb.com/v2.1.5/reference/출금-허용-주소-리스트-조회', slug: '출금-허용-주소-리스트-조회', file: '출금-허용-주소-리스트-조회.md' },
  { url: 'https://apidocs.bithumb.com/v2.1.5/reference/입금-리스트-조회', slug: '입금-리스트-조회', file: '코인-입금-리스트-조회.md' },
  { url: 'https://apidocs.bithumb.com/v2.1.5/reference/원화-입금-리스트-조회', slug: '원화-입금-리스트-조회', file: '원화-입금-리스트-조회.md' },
  { url: 'https://apidocs.bithumb.com/v2.1.5/reference/개별-입금-조회', slug: '개별-입금-조회', file: '개별-입금-조회.md' },
  { url: 'https://apidocs.bithumb.com/v2.1.5/reference/입금-주소-생성-요청', slug: '입금-주소-생성-요청', file: '입금-주소-생성-요청.md' },
  { url: 'https://apidocs.bithumb.com/v2.1.5/reference/전체-입금-주소-조회', slug: '전체-입금-주소-조회', file: '전체-입금-주소-조회.md' },
  { url: 'https://apidocs.bithumb.com/v2.1.5/reference/개별-입금-주소-조회', slug: '개별-입금-주소-조회', file: '개별-입금-주소-조회.md' },
  { url: 'https://apidocs.bithumb.com/v2.1.5/reference/원화-입금하기', slug: '원화-입금하기', file: '원화-입금하기.md' },
  { url: 'https://apidocs.bithumb.com/v2.1.5/reference/입출금-현황', slug: '입출금-현황', file: '입출금-현황.md' },
  { url: 'https://apidocs.bithumb.com/v2.1.5/reference/api-키-리스트-조회', slug: 'api-키-리스트-조회', file: 'API-키-리스트-조회.md' },
  { url: 'https://apidocs.bithumb.com/v2.1.5/reference/메이저-자산-조회담보및거래가능', slug: '메이저-자산-조회담보및거래가능', file: '메이저-자산-조회(담보및거래가능).md' },
  { url: 'https://apidocs.bithumb.com/v2.1.5/reference/상환레벨-조회', slug: '상환레벨-조회', file: '상환레벨-조회.md' }
];

function generateMarkdown(data) {
  let md = `# ${data.title}\n\n`;
  md += `> **원본 URL**: https://apidocs.bithumb.com/v2.1.5/reference/${data.urlSlug}\n\n`;

  if (data.method) {
    md += `## HTTP Method\n\n**${data.method}**\n\n`;
  }

  if (data.endpoint) {
    md += `## Endpoint\n\n\`\`\`\n${data.endpoint}\n\`\`\`\n\n`;
  }

  if (data.description) {
    md += `## 설명\n\n${data.description}\n\n`;
  }

  // Headers
  if (data.headers && data.headers.length > 0) {
    md += `## Headers\n\n`;
    md += `| 파라미터 | 타입 | 필수여부 | 설명 |\n`;
    md += `|---------|------|----------|------|\n`;
    data.headers.forEach(h => {
      md += `| ${h.name} | ${h.type} | ${h.required || ''} | ${h.description} |\n`;
    });
    md += '\n';
  }

  // Query Parameters
  if (data.queryParams && data.queryParams.length > 0) {
    md += `## Query Parameters\n\n`;
    md += `| 파라미터 | 타입 | 필수여부 | 설명 |\n`;
    md += `|---------|------|----------|------|\n`;
    data.queryParams.forEach(p => {
      md += `| ${p.name} | ${p.type} | ${p.required || ''} | ${p.description} |\n`;
    });
    md += '\n';
  }

  // Body Parameters
  if (data.bodyParams && data.bodyParams.length > 0) {
    md += `## Body Parameters\n\n`;
    md += `| 파라미터 | 타입 | 필수여부 | 설명 |\n`;
    md += `|---------|------|----------|------|\n`;
    data.bodyParams.forEach(p => {
      md += `| ${p.name} | ${p.type} | ${p.required || ''} | ${p.description} |\n`;
    });
    md += '\n';
  }

  // Response
  if (data.responseFields && data.responseFields.length > 0) {
    md += `## Response\n\n`;
    md += `| 필드 | 설명 | 타입 |\n`;
    md += `|------|------|------|\n`;
    data.responseFields.forEach(f => {
      md += `| ${f.field} | ${f.description} | ${f.type} |\n`;
    });
    md += '\n';
  }

  // Example Code
  if (data.exampleCode) {
    md += `## 예시 코드\n\n\`\`\`javascript\n${data.exampleCode}\n\`\`\`\n\n`;
  }

  // Response Example
  if (data.responseExample) {
    md += `## 응답 예시\n\n\`\`\`json\n${data.responseExample}\n\`\`\`\n\n`;
  }

  return md;
}

async function extractApiDoc(page, url, urlSlug) {
  try {
    await page.goto(url, { waitUntil: 'domcontentloaded', timeout: 30000 });
    await page.waitForTimeout(3000);

    const content = await page.evaluate(() => {
      const article = document.querySelector('article');
      if (!article) return null;

      const h1 = article.querySelector('h1');
      const title = h1 ? h1.textContent.trim() : '';

      const methodEl = article.querySelector('.rm-Methods');
      const method = methodEl ? methodEl.textContent.trim().toUpperCase() : '';

      const urlEl = article.querySelector('.rm-Spec-EndpointUrl a, .rm-Spec-EndpointUrl span');
      const endpoint = urlEl ? urlEl.textContent.trim() : '';

      const paragraphs = Array.from(article.querySelectorAll('p'));
      let description = '';
      for (const p of paragraphs) {
        const text = p.textContent.trim();
        if (text.length > 10 && !text.includes('예시코드') && !text.includes('JavaScript')) {
          description = text;
          break;
        }
      }

      const table = article.querySelector('table');
      let responseFields = [];
      if (table) {
        const rows = Array.from(table.querySelectorAll('tbody tr'));
        responseFields = rows.map(row => {
          const cells = Array.from(row.querySelectorAll('td'));
          return {
            field: cells[0]?.textContent?.trim() || '',
            description: cells[1]?.textContent?.trim() || '',
            type: cells[2]?.textContent?.trim() || ''
          };
        });
      }

      const headers = [];
      const headerSections = article.querySelectorAll('.rm-Param-section');
      headerSections.forEach(section => {
        const name = section.querySelector('.rm-Param-name')?.textContent?.trim() || '';
        const type = section.querySelector('.rm-Param-type')?.textContent?.trim() || '';
        const desc = section.querySelector('p')?.textContent?.trim() || '';
        const requiredEl = section.querySelector('.rm-Param-required');
        const required = requiredEl ? requiredEl.textContent.trim() : '';
        if (name && !name.includes('Language')) {
          headers.push({ name, type, description: desc, required });
        }
      });

      const queryParams = [];
      const querySections = article.querySelectorAll('[class*="query"]');
      querySections.forEach(section => {
        if (!section.closest('.rm-Param-section')) {
          const name = section.querySelector('.rm-Param-name')?.textContent?.trim() || '';
          const type = section.querySelector('.rm-Param-type')?.textContent?.trim() || '';
          const desc = section.querySelector('p')?.textContent?.trim() || '';
          const requiredEl = section.querySelector('.rm-Param-required');
          const required = requiredEl ? requiredEl.textContent.trim() : '';
          if (name) {
            queryParams.push({ name, type, description: desc, required });
          }
        }
      });

      const codeLines = Array.from(article.querySelectorAll('[data-line-number]'));
      let codeExample = '';
      if (codeLines.length > 0) {
        codeExample = codeLines.map(line => line.textContent.trim()).join('\n');
      }

      return {
        title,
        method,
        endpoint,
        description,
        responseFields,
        headers,
        queryParams,
        bodyParams: [],
        codeExample,
        responseExample: ''
      };
    });

    content.urlSlug = urlSlug;
    return content;
  } catch (error) {
    console.error(`Error extracting ${url}:`, error.message);
    return null;
  }
}

async function main() {
  const browser = await chromium.launch({ headless: true });
  const page = await browser.newPage();
  const outputDir = path.join(__dirname, '../docs/api/private');

  const results = [];

  for (let i = 0; i < API_LIST.length; i++) {
    const api = API_LIST[i];
    console.log(`[${i + 1}/${API_LIST.length}] Extracting: ${api.file}`);

    const apiData = await extractApiDoc(page, api.url, api.slug);
    if (apiData) {
      const markdown = generateMarkdown(apiData);
      const filePath = path.join(outputDir, api.file);

      await fs.writeFile(filePath, markdown, 'utf-8');
      console.log(`  ✓ Saved: ${filePath}`);

      results.push({
        file: api.file,
        title: apiData.title,
        endpoint: apiData.endpoint,
        method: apiData.method
      });
    } else {
      console.log(`  ✗ Failed to extract: ${api.url}`);
    }

    await page.waitForTimeout(1000);
  }

  await browser.close();

  console.log('\n=== Summary ===');
  console.log(`Total APIs: ${API_LIST.length}`);
  console.log(`Successfully extracted: ${results.length}`);
  console.log(`Failed: ${API_LIST.length - results.length}`);

  return results;
}

main().catch(console.error);
