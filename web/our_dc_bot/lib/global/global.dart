import 'package:riverpod_annotation/riverpod_annotation.dart';

import 'package:our_dc_bot/api/api.dart';
import 'package:our_dc_bot/api/model.dart';

part 'global.g.dart';

final Provider<UIAPIHandler> apiHandlerProvider = Provider((ref) {
  const apiBaseURL = String.fromEnvironment('API_BASE_URL', defaultValue: '');

  final handler = UIAPIHandler(Uri.parse(apiBaseURL));
  return handler;
});

class AuthState {
  UserInformation user = UserInformation((b) => b
    ..id = 0
    ..discordId = ''
    ..guildId = ''
    ..name = ''
    ..avatarUrl = '');

  bool isLogin = false;
}

@riverpod
class AuthStateNotifier extends _$AuthStateNotifier {
  AuthStateNotifier() : super();

  @override
  AuthState build() {
    return AuthState();
  }

  void updateLoginStatus(UserInformation user) {
    state.user = user;
    state.isLogin = true;
  }

  void logout() {
    state = AuthState();
  }
}
