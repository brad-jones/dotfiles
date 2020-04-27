import 'dart:io';
import 'package:ini/ini.dart';
import 'package:dexeca/dexeca.dart';
import 'package:dexecve/dexecve.dart';
import 'package:scripts/src/dir.dart';
import 'package:scripts/src/guest.dart';

final originalAwsVaultExe = normalizeDir(Platform.isLinux ? '~/.linuxbrew/bin/aws-vault' : '~/scoop/apps/aws-vault/current/aws-vault.exe');

Future<Config> readAwsConfig([String path = '~/.aws/config']) async {
  var lines = await File(normalizeDir(path)).readAsLines();

  var k = 0;
  for (var v in lines) {
    if (v.contains('[default]')) {
      lines[k] = v.replaceFirst('[default]', '[not-default]');
    }
    k++;
  }

  return Config.fromStrings(lines);
}

Future<String> getMfaSerial(String profile) async {
  if (profile?.isEmpty ?? true) throw Exception('empty profile');
  var cfg = await readAwsConfig();
  var serial = cfg.get('profile ${profile}', 'mfa_serial');
  if (serial?.isEmpty ?? true) {
    return await getMfaSerial(
      cfg.get('profile ${profile}', 'source_profile'),
    );
  }
  return serial.replaceAll('"', '');
}

Future<String> getMfaProvider(String profile) async {
  if (profile?.isEmpty ?? true) throw Exception('empty profile');
  var cfg = await readAwsConfig();
  var provider = cfg.get('profile ${profile}', 'mfa_token_provider');
  if (provider?.isEmpty ?? true) {
    return await getMfaProvider(
      cfg.get('profile ${profile}', 'source_profile'),
    );
  }
  return provider.replaceAll('"', '');
}

Future<String> runMfaProvider(String provider) async {
  var args = provider.split(' ');
  var proc = await dexeca(args[0], args.sublist(1), inheritStdio: false);
  return proc.stdout.trim();
}

Future<void> main(List<String> argv) async {
  var action = argv.isNotEmpty ? argv[0] : '';

  switch (action) {
    case 'exec':
    case 'login':
    case 'rotate':
      var args = <String>[action];
      var env = <String, String>{};

      if (argv.length >= 2) {
        if (!argv[1].startsWith('-')) {
          env['AWS_MFA_SERIAL'] = await getMfaSerial(argv[1]);
          args.addAll([
            '--mfa-token',
            await runMfaProvider(await getMfaProvider(argv[1])),
          ]);

          // install https://addons.mozilla.org/en-US/firefox/addon/open-url-in-container/ to make this work
          if (action == 'login') {
            args.add('-s');
            args.addAll(argv.sublist(1));
            var res = await dexeca(
              originalAwsVaultExe,
              args,
              inheritStdio: false,
              environment: env,
            );
            await execOnHostIfGuest('firefox', [
              'ext+container:name=${argv[1]}&url=${Uri.encodeQueryComponent(res.stdout)}'
            ]);
          }
        }
      }

      args.addAll(argv.sublist(1));
      dexecve(originalAwsVaultExe, args, environment: env);
      break;

    default:
      dexecve(originalAwsVaultExe, argv);
  }
}
