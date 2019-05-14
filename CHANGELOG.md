# Changelog

<a name="1.4.0"></a>
## [1.4.0] - 2019-05-15
### Bug Fixes
- failed to parse `|` character in tag parameters
- failed to parse multibyte string in tag parameters

### Chore
- some fix
- update comments
- update README
- update README
- add comment and remove unnecessary code

### Code Refactoring
- unexport unnecessary exported struct's field
- unexported if not neccesary
- make an internal error clear
- remove unreachable code

### Performance Improvements
- improve performance in failure

### Tests
- add tests for validator options
- add struct complex failure benchmark
- add tag parser test
- increase test coverage

### BREAKING CHANGE

using Apply method instead of SetXXX methods in Validator


<a name="1.3.0"></a>
## [1.3.0] - 2019-05-12
### Code Refactoring
- made behavior simple
- using parseXXX functions

### BREAKING CHANGE

using utf8.RuneCountInString by len, min and max

removed `zero` from default tags (can use `empty` instead)


<a name="1.2.0"></a>
## [1.2.0] - 2019-05-11
### Performance Improvements
- now 30% faster than before

### Tests
- Add complex validation benchmark

### Pull Requests
- Merge pull request [#7](https://github.com/utahta/go-validator/issues/7) from utahta/improve-performance


<a name="1.1.0"></a>
## [1.1.0] - 2019-04-18
### Features
- add empty tag


<a name="1.0.0"></a>
## [1.0.0] - 2019-03-16
### Bug Fixes
- order of adapter applied

### Features
- add go.mod

### BREAKING CHANGE

adapter feature


<a name="0.0.7"></a>
## [0.0.7] - 2019-01-21
### Bug Fixes
- break backwards compatibility

### Pull Requests
- Merge pull request [#6](https://github.com/utahta/go-validator/issues/6) from utahta/fix-backwards-compatibility


<a name="0.0.6"></a>
## [0.0.6] - 2019-01-21
### Bug Fixes
- tag parser

### Pull Requests
- Merge pull request [#5](https://github.com/utahta/go-validator/issues/5) from utahta/fix-optional-in-array


<a name="0.0.5"></a>
## [0.0.5] - 2019-01-20
### Code Refactoring
- tests
- remove unnecessary code
- fix to optional tag behavior
- Fix to error handling at validate

### Pull Requests
- Merge pull request [#4](https://github.com/utahta/go-validator/issues/4) from utahta/refactor-tags


<a name="0.0.4"></a>
## [0.0.4] - 2018-11-04
### Features
- Export default validate functions

### Pull Requests
- Merge pull request [#3](https://github.com/utahta/go-validator/issues/3) from utahta/export-validate-func


<a name="0.0.3"></a>
## [0.0.3] - 2018-10-22
### Chore
- Update .chglog/config.yml

### Code Refactoring
- Remove unreachable code

### Features
- Pass context to each validate function

### Pull Requests
- Merge pull request [#2](https://github.com/utahta/go-validator/issues/2) from utahta/context


<a name="0.0.2"></a>
## [0.0.2] - 2018-10-21
### Chore
- Add codecov.yml
- Add git-chglog

### Features
- Add SetAdapters function to Validator

### Tests
- Add tests for Validator
- Add tests for validate functions
- Add tests for parser and cacher
- Add tests for Field struct
- Add tests for Error struct

### Pull Requests
- Merge pull request [#1](https://github.com/utahta/go-validator/issues/1) from utahta/increase-coverage

### BREAKING CHANGE

Removed SetAdapter function from Validator.


<a name="0.0.1"></a>
## 0.0.1 - 2018-10-08
[1.4.0]: https://github.com/utahta/go-validator/compare/1.3.0...1.4.0
[1.3.0]: https://github.com/utahta/go-validator/compare/1.2.0...1.3.0
[1.2.0]: https://github.com/utahta/go-validator/compare/1.1.0...1.2.0
[1.1.0]: https://github.com/utahta/go-validator/compare/1.0.0...1.1.0
[1.0.0]: https://github.com/utahta/go-validator/compare/0.0.7...1.0.0
[0.0.7]: https://github.com/utahta/go-validator/compare/0.0.6...0.0.7
[0.0.6]: https://github.com/utahta/go-validator/compare/0.0.5...0.0.6
[0.0.5]: https://github.com/utahta/go-validator/compare/0.0.4...0.0.5
[0.0.4]: https://github.com/utahta/go-validator/compare/0.0.3...0.0.4
[0.0.3]: https://github.com/utahta/go-validator/compare/0.0.2...0.0.3
[0.0.2]: https://github.com/utahta/go-validator/compare/0.0.1...0.0.2
