import 'package:flutter/material.dart';

import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:our_dc_bot/api/api.dart';
import 'package:our_dc_bot/api/exception.dart';
import 'package:our_dc_bot/api/model.dart';
import 'package:our_dc_bot/routers/enum.dart';
import 'package:our_dc_bot/global/global.dart';

class DashboardLayout extends ConsumerWidget {
  const DashboardLayout({super.key, required this.navigationShellPage});

  final StatefulNavigationShell navigationShellPage;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return Scaffold(
      appBar: AppBar(
        backgroundColor: Theme.of(context).colorScheme.inversePrimary,
        title: const Text('Our Discord Bot'),
      ),
      drawer: DashboardNavigationDrawer(navigationShell: navigationShellPage),
      body: DashboardLayoutPage(page: navigationShellPage),
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
      return page;
    }

    return FutureBuilder<UserInformation>(
      future: _getSelfInformation(context, ref),
      builder: (context, snapshot) {
        if (snapshot.connectionState == ConnectionState.done) {
          if (snapshot.hasError) {
            if (snapshot.error is UnauthorizedException) {
              ref.read(authStateNotifierProvider.notifier).logout();
              context.goNamed(RouterName.signIn.path);
              return const SizedBox();
            }

            return const FailedInitPage();
          }

          ref
              .read(authStateNotifierProvider.notifier)
              .updateLoginStatus(snapshot.data!);
          return page;
        }

        return const Center(
          child: CircularProgressIndicator(),
        );
      },
    );
  }
}

class DashboardNavigationDrawer extends StatelessWidget {
  const DashboardNavigationDrawer({super.key, required this.navigationShell});

  final StatefulNavigationShell navigationShell;

  @override
  Widget build(BuildContext context) {
    return NavigationDrawer(
      selectedIndex: navigationShell.currentIndex,
      onDestinationSelected: (selectedIndex) {
        navigationShell.goBranch(
          selectedIndex,
          initialLocation: selectedIndex == navigationShell.currentIndex,
        );
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

List<NavigationDestinationContent> destinations =
    <NavigationDestinationContent>[
  NavigationDestinationContent(RouterName.stickerManager.name,
      const Icon(Icons.settings_outlined), const Icon(Icons.settings)),
  NavigationDestinationContent(RouterName.myInfo.name,
      const Icon(Icons.info_outlined), const Icon(Icons.info)),
];
