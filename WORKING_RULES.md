# WORKING_RULES

## Commit

- 格式:
```
type(scope): subject
```

- type:
    - feat: 對使用者有價值的新功能或可見變更  
    >[!NOTICE]
    > 如果是單純畫面排版，沒有新的資料顯示或隱藏等等，會使用style
    - docs: 文檔修改或新增
    - test: 測試
    - chore: 雜項修改, 通常與構建過程或工具相關
    - ci: CI/CD, 流水線相關
    - refactor: 重構代碼
    - style: 代碼格式化、調整UI畫面或新增註解
    - fix: 修復錯誤


- scope: 可加可不加，但如果要明確標示還是建議加一下，這邊統一規範scope要是適用區塊而不是子標
    - erogs: 批評空間相關模組
    - vndb: VNDB相關模組  
    ...

- subject: 原則上全小寫英文，遇到指令名稱可使用中文  
Ex. `feat: show game release date in '查詢創作者' command`

## Branch

- 格式:
```
type/subject
```

>[!NOTICE]
> 目前預計只會使用兩層，不要開到三層讓分支複雜化

- type:
    - feature: 新的功能
    - enhancement: 改善舊有功能
    - bugfix: 修復錯誤
    - hotfix: 緊急修復錯誤 **允許Bypass**
    - release: 發布版本
    - test: 測試分支(不會合到main，專門記錄用)
    - chore: 雜項修改, 通常與構建過程或工具相關
    - refactor: 內部重構程式碼(跟使用者無關)
    - ci: CI/CD, 流水線相關 **允許Bypass**
    - style: 代碼格式化、調整UI畫面或新增註解
    - docs: 純粹修改文檔(如果是更新紀錄允許同release一起發布)

## PR

- 目前無特別規定，只要注意審閱者、指派者以及Label要加就好

## Merge & Rebase

- 只允許全部的commit在兩天以內，並且該分支無其他rebase紀錄的分支可使用rebase，其餘都使用merge保障完整更新紀錄

