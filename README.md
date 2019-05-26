# API Test

API-Test is a concurrent runner for API integration tests.

## Concepts

An API integration test is defined by a set of [scenarios](https://github.com/efritz/api-test/blob/master/docs/scenarios.md). Each scenario contains either an ordered or unordered sequence of [tests](https://github.com/efritz/api-test/blob/master/docs/tests.md). A scenario can depend on the successful completion of zero or more other scenarios. This allows one scenario to prepare data or state in the API being tested necessary for later scenarios to run.

## Installation

Simply run `go install github.com/efritz/api-test`.

## Usage

The following command line flags are applicable for all IJ commands.

| Name             | Short Flag | Description |
| -----------------| ---------- | ----------- |
| config           | f          | The path to the config file. If not supplied, `api-test.yaml` and `api-test.yml` are attempted in the current directory. |
| force-sequential |            | Disable parallel execution. This will enforce that only one test is active at a time. The order that tests run are dependent only on scenario and test dependencies and may not consistent between runs. |
| junit            | j          | The target file to write the JUnit report file. |
| no-color         |            | Disable colorized output. |

## License

Copyright (c) 2019 Eric Fritz

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
