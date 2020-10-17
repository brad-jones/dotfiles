import 'dart:io';
import 'package:dexeca/dexeca.dart';
import 'package:http/http.dart' as http;

Future<String> currentChezmoiVersion() async {}

Future<String> latestChezmoiVersion() async {
  var client = http.Client();
  var request = http.Request(
      'GET', Uri.parse('https://github.com/twpayne/chezmoi/releases/latest'))
    ..followRedirects = false;
  var response = await client.send(request);
  client.close();
  return response.headers['location'].split('/').last.replaceFirst('v', '');
}

Future<void> installChezmoiLinux(String version) async {}

Future<void> installChezmoiWindows(String version) async {}

Future<void> main(List<String> argv) async {
  var chezmoiV = await latestChezmoiVersion();
  var response = await http.get(
      'https://github.com/twpayne/chezmoi/releases/download/v${chezmoiV}/chezmoi_${chezmoiV}_windows_amd64.zip');
  await File('').writeAsBytes(response.bodyBytes);

  //dexeca('chezmoi', ['apply', '--debug']);
}

/*
chezmoiV="$(wget https://github.com/twpayne/chezmoi/releases/latest -O /dev/null 2>&1 | grep Location: | sed -r 's~^.*tag/v(.*?) \[.*~\1~g')";
sudo dnf install -y https://github.com/twpayne/chezmoi/releases/download/v$chezmoiV/chezmoi-$chezmoiV-x86_64.rpm;

# Install chezmoi
RmIfExists -Path $env:TEMP\chezmoi;
RmIfExists -Path $env:TEMP\chezmoi.zip;
RmIfExists -Path $env:USERPROFILE\.local\bin\chezmoi.exe;
$ErrorActionPreference = 'continue';
$chezmoiV = wget https://github.com/twpayne/chezmoi/releases/latest -O /dev/null 2>&1 | grep Location: | sed -r 's~^.*tag/v(.*?) \[.*~\1~g';
$ErrorActionPreference = 'stop';
Exec -ScriptBlock { wget https://github.com/twpayne/chezmoi/releases/download/v${chezmoiV}/chezmoi_${chezmoiV}_windows_amd64.zip -O $env:TEMP\chezmoi.zip; }
Exec -ScriptBlock { 7z x $env:TEMP\chezmoi.zip "-o${env:TEMP}\chezmoi"; }
New-Item -ItemType Directory -Force -Path $env:USERPROFILE\.local\bin;
Copy-Item -Path $env:TEMP\chezmoi\chezmoi.exe -Destination $env:USERPROFILE\.local\bin\chezmoi.exe;
RmIfExists -Path $env:TEMP\chezmoi.zip; RmIfExists -Path $env:TEMP\chezmoi;
if (!($env:PATH -like "*$env:UserProfile\.local\bin*")) {
	$env:PATH += ";$env:UserProfile\.local\bin";
	[Environment]::SetEnvironmentVariable("PATH", $env:PATH, "User");
}
*/
