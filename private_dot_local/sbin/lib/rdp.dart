import 'dart:io';
import 'dart:convert';
import 'package:args/args.dart';
import 'package:dexeca/dexeca.dart';
import 'package:path/path.dart' as p;
import 'package:scripts/src/guest.dart';
import 'package:prompts/prompts.dart' as prompts;

void displayHelp(ArgParser parser, {String errorMessage, int exitCode = 0}) {
  if (errorMessage?.isNotEmpty ?? false) {
    stderr.writeln(errorMessage);
    stderr.writeln();
  }
  stdout.writeln('Usage: rdp [OPTIONS...]');
  stdout.writeln();
  stdout.writeln('Options:');
  stdout.writeln(parser.usage);
  exit(exitCode);
}

bool isIpAddress(String value) {
  return RegExp(r'^(?!0)(?!.*\.$)((1?\d?\d|25[0-5]|2[0-4]\d)(\.|$)){4}$')
      .hasMatch(value);
}

Future<void> login(
  String host,
  String username,
  String password, {
  String gateway,
}) async {
  var exitCode = 0;

  var addArg = '/add:"${host}"';
  if (isIpAddress(host)) {
    addArg = '/generic:"TERMSRV/${host}"';
  }
  await runOnHostIfGuest('cmdkey', [
    addArg,
    '/user:"${username}"',
    '/pass:"${password}"',
  ]);

  if (gateway?.isNotEmpty ?? false) {
    var addArg = '/add:"${gateway}"';
    if (isIpAddress(host)) {
      addArg = '/generic:"TERMSRV/${gateway}"';
    }
    await runOnHostIfGuest('cmdkey', [
      addArg,
      '/user:"${username}"',
      '/pass:"${password}"',
    ]);
  }

  try {
    await runOnHostIfGuest('mstsc.exe', ['/v:"${host}"']);
  } catch (e, st) {
    print(e);
    print(st);
    exitCode = -1;
  } finally {
    await Future.delayed(Duration(seconds: 1));

    if (isIpAddress(host)) {
      await runOnHostIfGuest('cmdkey', [
        '/delete:"TERMSRV/${host}"',
      ]);
    } else {
      await runOnHostIfGuest('cmdkey', [
        '/delete:"${host}"',
      ]);
    }

    if (gateway?.isNotEmpty ?? false) {
      if (isIpAddress(gateway)) {
        await runOnHostIfGuest('cmdkey', [
          '/delete:"TERMSRV/${gateway}"',
        ]);
      } else {
        await runOnHostIfGuest('cmdkey', [
          '/delete:"${gateway}"',
        ]);
      }
    }
  }

  exit(exitCode);
}

Future<ProcessResult> aws(List<String> args, [String profile]) async {
  var awsArgs = <String>[];
  if (profile?.isNotEmpty ?? false) {
    awsArgs.addAll(['--profile', profile]);
  }
  awsArgs.addAll(args);
  return await dexeca('aws', awsArgs, inheritStdio: false);
}

Future<Map<String, dynamic>> getEc2InstanceMetaData(
  String idOrName, [
  String profile,
]) async {
  if (RegExp(r'^i-.*$').hasMatch(idOrName)) {
    return jsonDecode((await aws([
      'ec2',
      'describe-instances',
      '--instance-ids',
      idOrName,
    ], profile))
        .stdout)['Reservations'][0]['Instances'][0];
  }

  return jsonDecode((await aws([
    'ec2',
    'describe-instances',
    '--filters',
    'Name=tag:Name,Values=${idOrName}',
  ], profile))
      .stdout)['Reservations'][0]['Instances'][0];
}

Future<String> getWindowsEc2Password(
  String instanceId,
  String key, [
  String profile,
]) async {
  String tmpDir;
  if (!await File(key).exists()) {
    tmpDir = (await Directory.systemTemp.createTemp()).path;
    var keyPath = p.join(tmpDir, 'key.pem');
    await File(keyPath).writeAsString(key);
    key = keyPath;
  }

  Map<String, dynamic> out;
  try {
    out = jsonDecode(
      (await aws([
        'ec2',
        'get-password-data',
        '--instance-id',
        instanceId,
        '--priv-launch-key',
        key,
      ], profile))
          .stdout,
    );
  } finally {
    if (tmpDir != null) {
      await Directory(tmpDir).delete(recursive: true);
    }
  }

  return out['PasswordData'] as String;
}

Future<void> main(List<String> argv) async {
  var parser = ArgParser();
  parser.addFlag(
    'help',
    help: 'displays this help screen',
    negatable: false,
  );
  parser.addOption(
    'host',
    abbr: 'h',
    help: 'The hostname or IP address of the server to connect to.',
    valueHelp: 'hostname-or-ip',
  );
  parser.addOption(
    'gateway',
    abbr: 'g',
    help: 'The RD Gateway to tunnel this connection through.',
    valueHelp: 'hostname-or-ip',
  );
  parser.addOption(
    'username',
    abbr: 'u',
    help: 'The username used to authenticate with the remote system.',
    valueHelp: 'DOMAIN\\User',
  );
  parser.addOption(
    'password',
    abbr: 'p',
    help: 'The password used to authenticate with the remote system.\n' +
        'If none provided a prompt will ask for it.',
  );
  parser.addOption(
    'ec2',
    help: 'An EC2 instance to connect to, use of this means the host,\n' +
        'username and password options will be ignored.',
    valueHelp: 'instance-id-or-name',
  );
  parser.addOption(
    'profile',
    help: 'The AWS profile as configured in ~/.aws/config to use to \n' +
        'access the ec2 instance. If not provided then this tool \n' +
        'assumes the environment is pre-authenticated.',
  );
  parser.addOption(
    'key',
    help: 'The SSH key to use to retrieve the EC2 password.',
    valueHelp: 'filepath-or-keybytes',
  );
  var parsedArgv = parser.parse(argv);

  if (parsedArgv.wasParsed('help') || argv.isEmpty) {
    displayHelp(parser);
  }

  if (parsedArgv.wasParsed('host') && parsedArgv.wasParsed('username')) {
    String password;
    if (!parsedArgv.wasParsed('password')) {
      password = prompts.get('Password', conceal: true);
    } else {
      password = parsedArgv['password'];
    }
    await login(
      parsedArgv['host'],
      parsedArgv['username'],
      password,
      gateway: parsedArgv.wasParsed('gateway') ? parsedArgv['gateway'] : null,
    );
  }

  if (parsedArgv.wasParsed('ec2') && parsedArgv.wasParsed('key')) {
    String profile;
    if (parsedArgv.wasParsed('profile')) {
      profile = parsedArgv['profile'];
    }

    var username = 'Administrator';
    if (parsedArgv.wasParsed('username')) {
      username = parsedArgv['username'];
    }

    var meta = await getEc2InstanceMetaData(parsedArgv['ec2'], profile);
    var pwd = await getWindowsEc2Password(
      meta['InstanceId'],
      parsedArgv['key'],
      profile,
    );
    await login(
      meta['PrivateIpAddress'],
      username,
      pwd,
      gateway: parsedArgv.wasParsed('gateway') ? parsedArgv['gateway'] : null,
    );
  }

  displayHelp(parser,
      errorMessage:
          'At minimum you must either supply the --host & --username options ' +
              'or the --ec2 & --key options!',
      exitCode: -1);
}
