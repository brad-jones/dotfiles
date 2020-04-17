import 'dart:io';
import 'package:path/path.dart' as p;

String normalizeDir(String path) {
  return p.normalize(path.replaceFirst('~', getHomeDir()));
}

String getHomeDir() {
  switch (Platform.operatingSystem) {
    case 'windows': return Platform.environment['UserProfile'];
    default: return Platform.environment['HOME'];
  }
}