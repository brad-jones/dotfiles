import 'dart:io';
import 'dart:convert';
import 'package:utf/utf.dart';
import 'package:dexeca/dexeca.dart';
import 'package:dexecve/dexecve.dart';

bool __runningAsGuestValue = false;
var _runningAsGuestChecked = false;
Future<bool> runningAsGuest() async {
  if (!_runningAsGuestChecked) {
    var hostsFile = File('/etc/hosts');
    if (await hostsFile.exists()) {
      if ((await hostsFile.readAsString()).contains('dev-server.wslhv.local')) {
        __runningAsGuestValue = true;
      }
    }
  }
  return __runningAsGuestValue;
}

Future<ProcessResult> runOnHostIfGuest(String exe, List<String> args) async {
  if (await runningAsGuest()) {
    return await dexeca('ssh', [
      '-o',
      'StrictHostKeyChecking=no',
      'dev-server.wslhv.local',
      exe,
      ...args,
    ]);
  }

  return await dexeca(exe, args);
}

Future<void> execOnHostIfGuest(String exe, List<String> args) async {
  if (await runningAsGuest()) {
    dexecve('ssh', [
      '-o',
      'StrictHostKeyChecking=no',
      'dev-server.wslhv.local',
      exe,
      ...args,
    ]);
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
