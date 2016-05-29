#!/usr/bin/env python
# encoding:utf-8
#
# Copyright 2015-2016 Yoshihiro Tanaka
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

__Author__ =  "Yoshihiro Tanaka <contact@cordea.jp>"
__date__   =  "2015-12-25"

import os, sys
import re

# @annotation
#   public
#   private
#   protected
# public
# private
# protected
# ...

with open(sys.argv[1]) as f:
    lines = f.readlines()

indent = re.compile(r'\ {4}')
deepIndent = re.compile(r'\ {8,}')
method = re.compile(r'\ {4}(?:public|private|protected)?[^(?:class)]?.*\(.*\)[ \t]*{')
endMethod = re.compile(r'\ {4}}')
var = re.compile(r'\ {4}(private|public|protected)?[\ \t]*[A-Za-z]+\w*(:?<\w+>)?[\ \t]+[A-Za-z]+\w*[\ \t]*=?[\ \t]*.*;$')
annotation = re.compile(r'\ {4}@([A-Za-z]+\w*)[ \t\w\(\)\.]*$')
cls = re.compile(r'(?:private|public|protected)?[ \t]*class[ \t]*\w+.*')
clsEnd = re.compile(r'(?:\ {4,})?.*{[\ \t]*')

isInMethod = False

variables = {}

sortList = ("public", "protected", "private", "z")

removeLines = []

for i in range(len(lines)):
    line = lines[i].rstrip()
    if deepIndent.match(line) != None:
        continue
    if endMethod.match(line) != None:
        isInMethod = False
        continue
    if isInMethod:
        continue
    if method.match(line) != None:
        isInMethod = True
        continue
    prevLine = lines[i-1].rstrip()
    m = var.match(line)
    if m:
        n = annotation.match(prevLine)
        optLine = None
        if not n:
            if not var.match(prevLine):
                optLine = prevLine
        mm = m.group(1) if m.group(1) else "z"
        mm = sortList.index(mm)
        nn = n.group(1) if n and n.group(1) else "z"
        prevLine = prevLine if n else None
        line = optLine + '\n' + line if optLine else line
        if nn in variables:
            if mm in variables[nn]:
                variables[nn][mm].append((prevLine, line))
            else:
                variables[nn][mm] = [(prevLine, line)]
        else:
            variables[nn] = {}
            variables[nn][mm] = [(prevLine, line)]
        removeLines.append(i)
        if nn != "z" or optLine:
            removeLines.append(i - 1)

mvs = ""
for key in sorted(variables.keys()):
    for inKey in sorted(variables[key].keys()):
        for value in sorted(variables[key][inKey]):
            if value[0]:
                mvs += value[0] + '\n'
            mvs += value[1] + '\n\n'
print removeLines

isClass = False
isInMethod = False
isOut = False
for i in range(len(lines)):
    line = lines[i].rstrip()
    if i in removeLines:
        continue

    if isClass:
        if not isOut:
            print
            print mvs
            isOut = True
        if not isInMethod:
            if line == '':
                continue
    if method.match(line):
        isInMethod = True
    if clsEnd.match(line):
        if cls.match(line) or cls.match(lines[i-1].rstrip()):
            isClass = True

    print line
