# Typescript Exporter
다국어 데이터를 사용할 수 있는 Typescript 패키지를 만들어 주는 exporter입니다.

## 메타데이터 파일 준비
`metadata.json`에서 `exporter_options`의 `typescript` 키 아래 옵션을 지정할 수 있습니다.
다음과 같은 정보가 필요합니다.
- `package_name`: 생성되는 라이브러리의 npm 패키지명

## 사용 방법 <span id="usage"></span>
생성된 패키지의 모듈 최상위에서는 아래와 같은 값들을 export합니다.
```typescript
export { Donggu } from "./donggu";
export { Version, RequiredLanguage, Language } from "./generated/dictionary";
```
- `Donggu`: 텍스트를 접근하기 위해 사용해야 하는 클래스입니다. `Donggu` 객체를 통해 생성된 텍스트에 접근합니다.
- `Version`: 라이브러리의 버전. `metadata.json`에 정의된 버전과 같습니다.
- `Language`: `metadata.json`에 정의된 `language` 값들의 union type입니다.
- `RequiredLanguage`: `metadata.json`에 정의된 `required_language` 값들의 union type입니다.

### Donggu Class
다국어 데이터는 Donggu 클래스를 통해 접근할 수 있습니다.
Donggu 생성자는 **사용자의 선호 언어를 반환하는 함수를**필요로 합니다.
`Language`의 배열을 반환해야 하며, 첫번째 언어부터 순서대로 우선순위를 가집니다.
예를 들어 `[en, ko, de]`를 반환했다면 영어, 한국어, 독일어 순으로 텍스트 사용을 시도합니다.

```ts
import { Donggu } from "@my-org/translation";

const donggu = new Donggu(() => {
    // 언어 정보를 가져오는 방법에 대한 예시
    const languages = getUserLanguagePreference();
    // 'en'이 RequiredLangauge에 있을 경우
    return [...langauges, 'en'];
});
```
사용하는 상황이 Express 서버라면 HTTP 요청을 이 함수에서 파싱할 것이고,
웹 애플리케이션이라면 브라우저의 설정을 파싱할 것입니다.

Node.js 서버에서 Donggu를 연동하는 방법은 [연동 가이드](#integration)에서 자세히 설명합니다.

함수가 반환하는 타입은 `[...Language[], RequiredLanguage]`이기 때문에, 항상 배열의 마지막에
모든 텍스트가 지원하는 언어를 지정해 주어야 합니다. 이 값은 기본 언어의 역할을 하며, 텍스트에 연결된
언어가 없어 텍스트 출력을 하지 못하는 오류를 막기 위한 안전장치입니다.

### 함수 사용
다국어 데이터들은 키만 `camelCase`로 변경되어 그대로 사용할 수 있습니다.
```json
{
    "screens.login_page.modal.success": {
        "ko": "로그인 성공"
    },
    "screens.login_page.ban_message": {
        "ko": "#{USER_NAME|string}님은 #{MONTH|int|02}월 #{DAY|int|02}일까지 제한된 사용자입니다."
    }
}
```
이런 데이터 항목들이 있다면 
```ts
donggu.screens.loginPage.modal.success()
// 로그인 성공
```
함수를 호출하여 텍스트를 가져올 수 있습니다. 언어는 Donggu 인스턴스를 만들때 넘겨주었던 선호 언어 판단 함수의 결과로 결정됩니다.

템플릿들은 함수의 매개변수로 넘겨주며, 템플릿 자료형에 지정한 타입을 그대로 사용합니다. 템플릿 키는 `camelCase`로 변환됩니다.
```ts
donggu.screens.loginPage.banMessage({userName: "홍길동", month: 8, day: 31 });
// 홍길동님은 08월 31일까지 제한된 사용자입니다.
```


## 연동 가이드<span id="integration"></span>

### Express / Nest.js와 연동

#### 1. HTTP 헤더의 `Accept-Language` 값을 파싱하여 사용자의 선호 언어를 추출합니다.
헤더 파싱은 직접 해도 괜찮지만,
여기서는 [accept-language-parser](https://www.npmjs.com/package/accept-language-parser) 패키지를 사용합니다.

#### 2. 가장 선호하는 언어로 텍스트를 생성할 수 있는 `Donggu` 인스턴스를 생성하고, Express의 request 객체에 저장합니다.

Express 애플리케이션을 준비하는 단계에서 아래와 같이 미들웨어를 추가합니다.
동구를 사용하는 미들웨어나 라우트 전에 이 미들웨어를 배치해야 합니다.
```ts
import { Donggu } from "@my-org/translation";
import { parse as parseAcceptLanguage } from "accept-language-parser";
// ...
const app = express();
// 다른 준비과정...
app.use((req, res, next) => {
    const parsedHeader = parseAcceptLanguage(req.headers['accept-language'] ?? 'en');
    const languages = parsedHeaders.map(item => item.code);

    request.DG = new Donggu(() => languages);
    next();
});
```

Express의 Request 객체에 속성을 추가한 것이기 때문에, 타입 에러가 나지 않으려면 `.d.ts` 파일을 만들어
Express의 타입 정의도 확장해 주어야 합니다.
```ts
import { Donggu } from "@my-org/translation";

declare global {
    namespace Express {
        interface Request {
            DG: Donggu;
        }
    }
}
```

#### 3. 텍스트를 생성해야 하는 라우트에서 `Donggu` 인스턴스를 가져와 사용합니다.
```ts
app.post("/login", (req, res) => {
    loginService.login(req).then((success) => {
        if (success) {
            res.status(200).json({ message: req.DG.login.success() });
        } else {
            res.status(401).json({ message: req.DG.login.unauthorized() });
        }
    });
});
```
### `cls-hooked`를 사용해 Express / Nest.js와 연동
[cls-hooked](https://www.npmjs.com/package/cls-hooked)를 사용하면 Express request에 직접 동구 인스턴스를 넣을 필요 없이
사용자의 언어 선호 순서를 저장, 전달할 수 있습니다.

cls-hooked에 대한 설명은 여기서 하지 않으나, cls-hooked의 작동 방식을 먼저 이해하고 이 방법을 사용하는 것을 권장합니다.

#### 0. `cls-hooked`에 동구의 namespace와 namespace를 이용하는 `Donggu` 인스턴스를 생성합니다.
별도의 파일에 아래와 같이 네임스페이스를 정의, 사용하는 기능을 작성합니다.
```ts
// donggu.ts
import { createNamespace } from "cls-hooked";
import { Language, RequiredLanguage } from "@my-org/translation";

type LanguagePreference = [...Language[], RequiredLanguage];

const NAMESPACE_KEY = "LANGUAGE_NS";
const ENTRY_KEY = "LANG_PREF";
const ns = createNamespace(NAMESPACE_KEY);

export function wrapInDongguCls(language: LanguagePreference, func: () => void) {
    ns.run(() => {
        ns.set(ENTRY_KEY, language);
        func();
    });
}

export function getDongguPreference(): LanguagePreference | undefined {
    return (ns.get(ENTRY_KEY) ?? undefined) as (LanguagePreference | undefined);
}

export const DG = new Donggu(() => {
    // 'en'이 RequiredLanguage에 포함되어 있지 않을 경우 다른 언어로 변경해야 합니다.
    return getDongguPreference() ?? ['en'];
});
```

1. HTTP 헤더의 `Accept-Language` 값을 파싱하여 사용자의 선호 언어를 추출합니다.
2. cls-hooked를 제대로 쓸 수 있도록 래퍼로 함수를 감싸고, 선호 언어를 **동구 namespace에** 저장합니다.
3. 텍스트를 생성해야 하는 라우트에서 `Donggu` 인스턴스를 가져와 사용합니다.

미들웨어는 이렇게 정의할 수 있습니다.
```ts
import { Donggu } from "@my-org/translation";
import { parse as parseAcceptLanguage } from "accept-language-parser";
import { wrapInDongguCls } from "../donggu";
// ...
const app = express();
// 다른 준비과정...
app.use((req, res, next) => {
    const parsedHeader = parseAcceptLanguage(req.headers['accept-language'] ?? 'en');
    const languages = parsedHeaders.map(item => item.code);

    wrapInDongguCls(languages, () => {
        next();
    });
});
```

텍스트를 생성해야 하는 라우트에서는 `donggu.ts`에서 만든 `Donggu` 인스턴스를 가져와 사용합니다.
```ts
import { DG } from "../donggu";
// ...
app.post("/login", (req, res) => {
    loginService.login(req).then((success) => {
        if (success) {
            res.status(200).json({ message: DG.login.success() });
        } else {
            res.status(401).json({ message: DG.login.unauthorized() });
        }
    });
});
```
