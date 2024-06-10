import 'package:flutter/material.dart';

import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:our_dc_bot/api/api.dart';
import 'package:our_dc_bot/api/exception.dart';
import 'package:our_dc_bot/api/model.dart';
import 'package:our_dc_bot/routers/enum.dart';
import 'package:our_dc_bot/global/global.dart';

class DashboardLayout extends ConsumerWidget {
  const DashboardLayout({super.key, required this.page});

  final Widget page;

  Future<UserInformation> _getSelfInformation(
      BuildContext context, WidgetRef ref) async {
    return ref
        .read<UIAPIHandler>(apiHandlerProvider)
        .call(context, (api) => api.getSelfInformation());
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    AuthState authState = ref.watch(authStateNotifierProvider);
    if (authState.isLogin) {
      return DashboardLayoutPage(page: page);
    }

    return FutureBuilder<UserInformation>(
      future: _getSelfInformation(context, ref),
      builder: (context, snapshot) {
        if (snapshot.connectionState == ConnectionState.done) {
          if (snapshot.hasError) {
            if (snapshot.error is UnauthorizedException) {
              ref.read(authStateNotifierProvider.notifier).logout();
              context.goNamed(RouterName.signIn.path);
              return const DashboardLayoutPage(page: Center());
            }

            return const DashboardLayoutPage(page: FailedInitPage());
          }

          ref
              .read(authStateNotifierProvider.notifier)
              .updateLoginStatus(snapshot.data!);
          return DashboardLayoutPage(page: page);
        }

        return const DashboardLayoutPage(
          page: Center(
            child: CircularProgressIndicator(),
          ),
        );
      },
    );
  }
}

class FailedInitPage extends StatelessWidget {
  const FailedInitPage({super.key});

  @override
  Widget build(BuildContext context) {
    return const Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: <Widget>[
          Text('Failed to initialize'),
          Text('Please try again later'),
        ],
      ),
    );
  }
}

class DashboardLayoutPage extends ConsumerWidget {
  const DashboardLayoutPage({
    super.key,
    required this.page,
  });

  final Widget page;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return Scaffold(
      appBar: AppBar(
        backgroundColor: Theme.of(context).colorScheme.inversePrimary,
        title: const Text('Our Discord Bot'),
      ),
      drawer: const DashboardNavigationDrawer(),
      body: page,
    );
  }
}

class DashboardNavigationDrawer extends ConsumerStatefulWidget {
  const DashboardNavigationDrawer({super.key});

  @override
  DashboardNavigationDrawerState createState() =>
      DashboardNavigationDrawerState();
}

class DashboardNavigationDrawerState
    extends ConsumerState<DashboardNavigationDrawer> {
  int _navigationSelectedIndex = 1;

  @override
  Widget build(BuildContext context) {
    return NavigationDrawer(
      selectedIndex: _navigationSelectedIndex,
      onDestinationSelected: (selectedScreen) {
        setState(() {
          _navigationSelectedIndex = selectedScreen;
          context.goNamed(destinations[selectedScreen].label);
        });
      },
      children: <Widget>[
        Padding(
          padding: const EdgeInsets.fromLTRB(28, 16, 16, 10),
          child: Text(
            '導航',
            style: Theme.of(context).textTheme.titleSmall,
          ),
        ),
        ...destinations.map((NavigationDestinationContent destination) {
          return NavigationDrawerDestination(
            label: Text(destination.label),
            icon: destination.icon,
            selectedIcon: destination.selectedIcon,
          );
        }),
      ],
    );
  }
}

class NavigationDestinationContent {
  const NavigationDestinationContent(this.label, this.icon, this.selectedIcon);

  final String label;
  final Widget icon;
  final Widget selectedIcon;
}

const List<NavigationDestinationContent> destinations =
    <NavigationDestinationContent>[
  NavigationDestinationContent(
      '貼圖管理', Icon(Icons.settings_outlined), Icon(Icons.settings)),
  NavigationDestinationContent(
      '我的資訊', Icon(Icons.info_outlined), Icon(Icons.info)),
];
