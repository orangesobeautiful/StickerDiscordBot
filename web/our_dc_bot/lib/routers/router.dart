import 'package:flutter/material.dart';

import 'package:go_router/go_router.dart';

import 'package:our_dc_bot/layouts/dashboard.dart';
import 'package:our_dc_bot/pages/my_info.dart';
import 'package:our_dc_bot/pages/signin.dart';
import 'package:our_dc_bot/routers/enum.dart';

final _rootNavigatorKey = GlobalKey<NavigatorState>();
final _shellNavigatorKey = GlobalKey<NavigatorState>();

RouterConfig<Object> newRouter() {
  return GoRouter(
    navigatorKey: _rootNavigatorKey,
    initialLocation: '/',
    routes: <RouteBase>[
      ShellRoute(
        navigatorKey: _shellNavigatorKey,
        builder: (BuildContext context, GoRouterState state, Widget child) {
          return DashboardLayout(
            page: child,
          );
        },
        routes: [
          GoRoute(
            parentNavigatorKey: _shellNavigatorKey,
            name: RouterName.myInfo.name,
            path: RouterName.myInfo.path,
            builder: (context, state) {
              return const MyInfoPage();
            },
          ),
          GoRoute(
            parentNavigatorKey: _shellNavigatorKey,
            name: RouterName.stickerManager.name,
            path: RouterName.stickerManager.path,
            builder: (context, state) {
              return const Center(
                child: Text('Sticker Manager'),
              );
            },
          ),
        ],
      ),
      GoRoute(
        name: RouterName.signIn.name,
        path: RouterName.signIn.path,
        parentNavigatorKey: _rootNavigatorKey,
        builder: (context, state) {
          return const SignInPage();
        },
      ),
    ],
  );
}
