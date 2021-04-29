# [1.7.0](https://github.com/brad-jones/dotfiles/compare/v1.6.1...v1.7.0) (2021-04-29)


### Bug Fixes

* deno dnf install command in the linux dotfile-updater script ([989f5e2](https://github.com/brad-jones/dotfiles/commit/989f5e257e3b7350485968f1415b3a9e9604d200))


### Features

* added docker into WSL via Podman ([f1b44db](https://github.com/brad-jones/dotfiles/commit/f1b44db7d2c826bce782bfff5dd9dad02e9eefec))

## [1.6.1](https://github.com/brad-jones/dotfiles/compare/v1.6.0...v1.6.1) (2021-04-29)


### Bug Fixes

* **deps:** bump github.com/brad-jones/goexec/v2 from 2.1.6 to 2.1.7 ([#11](https://github.com/brad-jones/dotfiles/issues/11)) ([899b439](https://github.com/brad-jones/dotfiles/commit/899b439fbf7f608ef749637962da4671307d3f09))

# [1.6.0](https://github.com/brad-jones/dotfiles/compare/v1.5.7...v1.6.0) (2021-04-26)


### Features

* added deno into the linux install ([547f729](https://github.com/brad-jones/dotfiles/commit/547f72921e4e37bd96fe15d82e4951737635b78c))

## [1.5.7](https://github.com/brad-jones/dotfiles/compare/v1.5.6...v1.5.7) (2021-04-26)


### Bug Fixes

* ensure the latest version of npm is installed also ([20ec960](https://github.com/brad-jones/dotfiles/commit/20ec96012a99add83592db2ab5f664619d17f0da))

## [1.5.6](https://github.com/brad-jones/dotfiles/compare/v1.5.5...v1.5.6) (2021-04-26)


### Bug Fixes

* ignore some of the dart files that should not have been comitted ([95bb41b](https://github.com/brad-jones/dotfiles/commit/95bb41be57e3db82706518974db1791841d603f8))
* **deps:** bump github.com/google/go-github/v35 from 35.0.0 to 35.1.0 ([#7](https://github.com/brad-jones/dotfiles/issues/7)) ([e60fadd](https://github.com/brad-jones/dotfiles/commit/e60faddd5944b867f222a8a8dc4b772374d98836))

## [1.5.5](https://github.com/brad-jones/dotfiles/compare/v1.5.4...v1.5.5) (2021-04-25)


### Bug Fixes

* cant grab the vault password from inside wsl ([2c5cdf2](https://github.com/brad-jones/dotfiles/commit/2c5cdf2874f33452518c1dce2ead1d256869407a))

## [1.5.4](https://github.com/brad-jones/dotfiles/compare/v1.5.3...v1.5.4) (2021-04-25)


### Bug Fixes

* retry the entire vault unlocking ([cede457](https://github.com/brad-jones/dotfiles/commit/cede45758172c5495f5806cafad51f4bb65660c6))

## [1.5.3](https://github.com/brad-jones/dotfiles/compare/v1.5.2...v1.5.3) (2021-04-25)


### Bug Fixes

* logon script installer didn't work with the updater ([15deaaf](https://github.com/brad-jones/dotfiles/commit/15deaaf218a5961cd645d7dc10539ad4a567fab1))

## [1.5.2](https://github.com/brad-jones/dotfiles/compare/v1.5.1...v1.5.2) (2021-04-25)


### Bug Fixes

* replace ping with normal http call ([a6e8d09](https://github.com/brad-jones/dotfiles/commit/a6e8d0974e7e73bbf35a10c09d95342124f9188c))

## [1.5.1](https://github.com/brad-jones/dotfiles/compare/v1.5.0...v1.5.1) (2021-04-25)


### Bug Fixes

* wait for internet access before trying to do anything ([7cb215e](https://github.com/brad-jones/dotfiles/commit/7cb215ef0c0b0ef08e7809503131130fc7b53899))

# [1.5.0](https://github.com/brad-jones/dotfiles/compare/v1.4.2...v1.5.0) (2021-04-25)


### Features

* **wsl:** good enough for now ([f9fdcc0](https://github.com/brad-jones/dotfiles/commit/f9fdcc02801af50b5383a33bcd530ade409329ef))

## [1.4.2](https://github.com/brad-jones/dotfiles/compare/v1.4.1...v1.4.2) (2021-04-25)


### Bug Fixes

* **deps:** bump github.com/AlecAivazis/survey/v2 from 2.2.9 to 2.2.12 ([#6](https://github.com/brad-jones/dotfiles/issues/6)) ([6f8d7fb](https://github.com/brad-jones/dotfiles/commit/6f8d7fb4649ce543d3f731c7c97621fb44ef8d23))

## [1.4.1](https://github.com/brad-jones/dotfiles/compare/v1.4.0...v1.4.1) (2021-04-21)


### Bug Fixes

* **selfupdate:** now still runs rest of process if no update needed ([a818717](https://github.com/brad-jones/dotfiles/commit/a818717aff9b59c79216cb366b104ef1a9858057))

# [1.4.0](https://github.com/brad-jones/dotfiles/compare/v1.3.1...v1.4.0) (2021-04-21)


### Bug Fixes

* **dartscripts:** remove some debugging code ([c172e28](https://github.com/brad-jones/dotfiles/commit/c172e289e5b3afb6735a69a038833f5ea8717d22))


### Features

* self update is now a thing ([0426e5a](https://github.com/brad-jones/dotfiles/commit/0426e5ab384b4e1334a8574c028f11ba4de1e386))

## [1.3.1](https://github.com/brad-jones/dotfiles/compare/v1.3.0...v1.3.1) (2021-04-21)


### Bug Fixes

* **gopass:** make sure we sync the vault after unlocking it ([6606ed5](https://github.com/brad-jones/dotfiles/commit/6606ed553513ca295b82090a55b979d824b99d9f))
* add some sleep time at the end of the program ([017de6e](https://github.com/brad-jones/dotfiles/commit/017de6eba3c6933ce0b97cf413ee56563ffb020d))
* **killprocbyname:** add some additional logging ([aa70cb3](https://github.com/brad-jones/dotfiles/commit/aa70cb3167528fa472521c0072b0e9518d478417))

# [1.3.0](https://github.com/brad-jones/dotfiles/compare/v1.2.1...v1.3.0) (2021-04-20)


### Bug Fixes

* **winsudo:** updated error txt to be correct ([20f574b](https://github.com/brad-jones/dotfiles/commit/20f574b1dbac842d33eb99c2fcb22220ed8082d6))


### Features

* added press enter to continue logic to error handler ([d2d7c58](https://github.com/brad-jones/dotfiles/commit/d2d7c58cf868e27d1a7cd59ad340b6274e1d52c6))

## [1.2.1](https://github.com/brad-jones/dotfiles/compare/v1.2.0...v1.2.1) (2021-04-20)


### Bug Fixes

* vault unlocking and updating, still not totally happy with this ([e285027](https://github.com/brad-jones/dotfiles/commit/e28502720904cea2625a5493c86f3cff7ec3f30d))

# [1.2.0](https://github.com/brad-jones/dotfiles/compare/v1.1.2...v1.2.0) (2021-04-19)


### Features

* executes the dotfiles bin directly on logon ([1594d5d](https://github.com/brad-jones/dotfiles/commit/1594d5dfa9acc4bb1d12aa180f2a4fbd6baccf3a))

## [1.1.2](https://github.com/brad-jones/dotfiles/compare/v1.1.1...v1.1.2) (2021-04-19)


### Bug Fixes

* downloaded URL checked against Notarized URL ([f754421](https://github.com/brad-jones/dotfiles/commit/f754421b2b94c8f9741ce377e2e3497f455a3df8))

## [1.1.1](https://github.com/brad-jones/dotfiles/compare/v1.1.0...v1.1.1) (2021-04-19)


### Bug Fixes

* powershell installer direct to null was incorrect ([5367a50](https://github.com/brad-jones/dotfiles/commit/5367a5000a56919e909903686437df3c8ecdeb59))

# [1.1.0](https://github.com/brad-jones/dotfiles/compare/v1.0.1...v1.1.0) (2021-04-19)


### Bug Fixes

* ran go mod tidy ([ca6b1dd](https://github.com/brad-jones/dotfiles/commit/ca6b1dde72079b9dbc08d07d163357a6bd40c323))
* use versioned github go module ([360c41d](https://github.com/brad-jones/dotfiles/commit/360c41d7b4647a7cf4d363b8e61a198a7c1b1b61))


### Features

* codenotary installers ([bfe2961](https://github.com/brad-jones/dotfiles/commit/bfe2961ce57c79efdedb5a9986d742acc950f116))
* installs run at logon script for windows ([dbec5c7](https://github.com/brad-jones/dotfiles/commit/dbec5c7776d6c27c738224aacf8d9caa22a169b6))

## [1.0.1](https://github.com/brad-jones/dotfiles/compare/v1.0.0...v1.0.1) (2021-04-14)


### Bug Fixes

* commit gopass native host binary that got ignored previously ([f31c051](https://github.com/brad-jones/dotfiles/commit/f31c051a646c6f8be58c25601324b0f866f8c98f))

# 1.0.0 (2021-04-14)


### Features

* initial commit ([3fa4be3](https://github.com/brad-jones/dotfiles/commit/3fa4be378045a42296dcdf54bab68eefe2d34372))
