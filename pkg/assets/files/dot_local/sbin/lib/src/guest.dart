import 'dart:io';
import 'dart:convert';
import 'package:utf/utf.dart';
import 'package:dexeca/dexeca.dart';
import 'package:dexecve/dexecve.dart';

Future<ProcessResult> runOnHostIfGuest(String exe, List<String> args) async {
  if (Platform.isLinux) {
    exe = '$exe.exe';
  }
  return await dexeca(exe, args);
}

Future<void> execOnHostIfGuest(String exe, List<String> args) async {
  if (Platform.isLinux) {
    exe = '$exe.exe';
  }
  dexecve(exe, args);
}

Future<ProcessResult> powershell(String script) {
  return runOnHostIfGuest(
    'powershell',
    [
      '-NoLogo',
      '-NoProfile',
      '-WindowStyle',
      'Hidden',
      '-Output',
      'XML',
      '-EncodedCommand',
      base64.encode(encodeUtf16le(script)),
    ],
  );
}
