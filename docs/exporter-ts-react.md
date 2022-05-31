# Typescript React Exporter
다국어 데이터를 사용할 수 있는 Typescript React 패키지를 만들어 주는 exporter입니다.

기본적인 기능은 [Typescript Exporter](exporter-ts.md) 동일하지만,
리액트에서 편리한 기능이 몇가지 추가됩니다.

- [사용 방법](#usage)
- [연동 가이드](#integration)

## 메타데이터 파일 준비
`metadata.json`에서 `exporter_options`의 `ts-react` 키 아래 옵션을 지정할 수 있습니다.
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
    const language = navigator.language.split("-")[0];
    // 'en'이 RequiredLangauge에 있을 경우
    return [...langauge, "en"];
});
```
위의 경우는 [Navigator.language](https://developer.mozilla.org/ko/docs/Web/API/Navigator/language)의 값을 확인해 바로 선호 언어로 사용합니다. 환경에 따라 언어를 가져오는 부분을 자유롭게 설정하면 됩니다. 예를 들어

- React Native라면 [Platform]()과 [NativeModules]()에서 나오는 휴대폰 설정을 사용
- 상태 관리자를 사용하는 경우라면 상태 관리자의 hook을 사용
- 웹이라면 쿠키나 localStorage의 값을 사용

Hook을 사용하는 간단한 예시는 [연동 가이드](#integration)에서 다룹니다.

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
```tsx
<div>
    {donggu.screens.loginPage.modal.success()}
</div>
// 로그인 성공
```
함수를 호출하여 텍스트를 가져올 수 있습니다. 언어는 Donggu 인스턴스를 만들때 넘겨주었던 선호 언어 판단 함수의 결과로 결정됩니다.

템플릿들은 함수의 매개변수로 넘겨주며, 템플릿 자료형에 지정한 타입을 그대로 사용합니다. 템플릿 키는 `camelCase`로 변환됩니다.
```tsx
<div>
    {donggu.screens.loginPage.banMessage({userName: "홍길동", month: 8, day: 31 })}
</div>
// 홍길동님은 08월 31일까지 제한된 사용자입니다.
```

### 줄바꿈 문자 대체
줄바꿈 문자를 지정하지 않는다면 줄바꿈에 `\n`이 사용됩니다. 다른 Element를 사용해야 하는 경우에는
함수에 `lineBreakElement` 옵션을 지정해 대체할 수 있습니다.
```json
{
    "hi": {"ko": "안녕하세요,\n좋은\n하루입니다."}
}
```
```tsx
<div>
    {donggu.hi({ lineBreakElement: <br /> })}
</div>
// 안녕하세요,<br />좋은<br />하루입니다.
```

Donggu 인스턴스의 `lineBreakElement` 속성을 지정해서 기본 줄바꿈 문자를 변경할 수 있습니다.
```ts
donggu.lineBreakElement = <br />;
```

```tsx
<div>
    {donggu.hi()}
</div>
// 안녕하세요,<br />좋은<br />하루입니다.
```

### 템플릿을 컴포넌트로 감싸기
텍스트 중의 일부 부분에 스타일을 적용해야 할 경우, 템플릿을 컴포넌트로 감싸 해결할 수 있습니다.
```json
{
    "mail": {"ko": "메일이 #{COUNT|int}건 있습니다."}
}
```
```tsx
<div>
    {donggu.mail({count: 10}, {wrappingElement: {name: OrangeBoldText}})}
</div>
// 메일이 <OrangeBoldText>10</OrangeBoldText>건 있습니다.
```

## 연동 가이드 <span id="integration"></span>

### 훅을 사용하지 않는 경우
브라우저 설정을 그대로 따라가거나, 휴대폰 설정을 그대로 따라가는 경우와 같이
Donggu 인스턴스를 하나만 만들어 싱글턴으로 사용하는 것이 편리합니다.

```ts
// donggu.ts
import { Donggu } from "@my-org/translation";

export const DG = new Donggu(() => {
    const language = navigator.language.split("-")[0];
    // 'en'이 RequiredLangauge에 있을 경우
    return [...langauge, "en"];
});
```

```tsx
import { DG } from "../donggu";

export const MyComponent = () => {
    return (
        <div>
            {donggu.hi()}
        </div>
    );
};
```

### 훅을 사용하는 경우
상태 관리자, Context API와 같이 전역에서 접근 가능한 데이터로 언어를 설정하고자 할때는
Donggu 인스턴스를 훅으로 감싸 사용하면 편리합니다.

```ts
// use-donggu.ts
import { Donggu } from "@my-org/translation";
import { useLanguageState } from "my-state-manager";

export const useDonggu = () => {
    const [language] = useLanguageState();
    // 'en'이 RequiredLangauge에 있을 경우
    return new Donggu(() => [language, "en"]);
};
```
