import 'package:flutter/material.dart';

import 'package:go_router/go_router.dart';

import 'package:our_dc_bot/layouts/dashboard.dart';
import 'package:our_dc_bot/pages/my_info.dart';
import 'package:our_dc_bot/pages/signin.dart';
import 'package:our_dc_bot/routers/enum.dart';
import 'package:our_dc_bot/pages/sticker_manager.dart';

final _rootNavigatorKey = GlobalKey<NavigatorState>();
final _shellNavigatorKey = GlobalKey<NavigatorState>();

RouterConfig<Object> newRouter() {
  return GoRouter(
    navigatorKey: _rootNavigatorKey,
    initialLocation: '/',
    routes: <RouteBase>[
      StatefulShellRoute.indexedStack(
        builder: (BuildContext _, GoRouterState __,
            StatefulNavigationShell navigationShell) {
          return DashboardLayout(
            navigationShellPage: navigationShell,
          );
        },
        branches: <StatefulShellBranch>[
          StatefulShellBranch(
            navigatorKey: _shellNavigatorKey,
            routes: <RouteBase>[
              GoRoute(
                name: RouterName.stickerManager.name,
                path: RouterName.stickerManager.path,
                builder: (context, state) {
                  return const StickerManagerPage();
                },
              ),
            ],
          ),
          StatefulShellBranch(
            routes: <RouteBase>[
              GoRoute(
                name: RouterName.myInfo.name,
                path: RouterName.myInfo.path,
                builder: (context, state) {
                  return const MyInfoPage();
                },
              ),
            ],
          )
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
