import 'dart:io';
import 'package:dexeca/dexeca.dart';
import 'package:dexecve/dexecve.dart';
import 'package:scripts/src/dir.dart';

final originalAwsExe = normalizeDir(Platform.isLinux ? '~/.local/bin/aws' : '~/scoop/apps/aws/current/bin/aws.exe');

String getAwsProfile(List<String> argv) {
  var k = 0;
  for (var v in argv) {
    if (v == '--profile') {
      return argv[++k];
    }
    k++;
  }

  if (Platform.environment.containsKey('AWS_PROFILE')) {
    return Platform.environment['AWS_PROFILE'];
  }

  if (Platform.environment.containsKey('AWS_DEFAULT_PROFILE')) {
    return Platform.environment['AWS_DEFAULT_PROFILE'];
  }

  return '';
}

List<String> getArgsWithoutProfile(List<String> argv) {
  var newArgv = <String>[];

  var skip = false;
  for (var v in argv) {
    if (skip) {
      skip = false;
      continue;
    }
    if (v == '--profile') {
      skip = true;
      continue;
    }
    newArgv.add(v);
  }

  return newArgv;
}

Future<Map<String, String>> getEnvFromAwsVault(String profile) async {
  var env = <String, String>{};
  if (profile?.isEmpty ?? true) {
    return env;
  }

  ProcessResult result;
  if (Platform.isWindows) {
    result = await dexeca(
      'aws-vault',
      ['exec', profile, '--', 'cmd.exe', '/C', 'SET'],
      inheritStdio: false,
    );
  } else {
    result = await dexeca(
      'aws-vault',
      ['exec', profile, '--', 'env'],
      inheritStdio: false,
    );
  }

  for (var line in result.stdout.replaceAll('\r\n', '\n').split('\n')) {
    if (line.contains('=') && line.startsWith('AWS_')) {
      var parts = line.split('=');
      env[parts[0]] = parts[1];
    }
  }

  return env;
}

Future<void> main(List<String> argv) async {
  dexecve(
    originalAwsExe,
    getArgsWithoutProfile(argv),
    environment: await getEnvFromAwsVault(getAwsProfile(argv)),
  );
}
