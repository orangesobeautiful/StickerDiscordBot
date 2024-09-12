enum RouterName {
  myInfo(name: 'MyInfo', path: '/'),
  signIn(name: 'SignIn', path: '/signin'),
  stickerManager(name: 'StickerManager', path: '/sticker-manager');

  final String name;
  final String path;

  const RouterName({required this.name, required this.path});
}
