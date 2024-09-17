/*
Copyright © 2023 Morgan Winbush <morgan@ynotlabsllc.com>
*/
package flutter

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

// flutterCmd represents the flutter command
var FlutterCmd = &cobra.Command{
	Use:   "flutter",
	Short: "A brief description of your command",
	Long:  `A Flutter CLI for making dope shit`,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Hugo",
	Long:  `All software has versions. This is Hugo's`,
	Run: func(cmd *cobra.Command, args []string) {
		flutterCLI("--version", "1.0.1")
	},
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new Flutter Project",
	Long:  `Create is for generating a new Flutter project with all Product Shop templated items`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Generateing new Flutter Project: %v \n", args[0])
		flutterCLI("create", args[0])
	},
}

var createScreen = &cobra.Command{
	Use:   "screen",
	Short: "Create A New Screen Component",
	Long:  `This creates a new flutter bloc module`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Generateing new Flutter Screen Components for: %v \n", args[0])

		createScreenFile(args[0], "")
	},
}

func createScreenFile(screenName string, destination string) {
	screenNameLower := strings.ToLower(screenName)
	if destination == "" {
		cmd := exec.Command("mkdir", screenNameLower)
		_, err := cmd.CombinedOutput()
		check(err)
	}

	var screenFile string = ""
	var cubitFile string = ""
	if destination != "" {
		screenFile = fmt.Sprintf("%s%s_screen.dart", destination, screenNameLower)
		cubitFile = fmt.Sprintf("%s%s_cubit.dart", destination, screenNameLower)
	} else {
		screenFile = fmt.Sprintf("%s/%s_screen.dart", screenNameLower, screenNameLower)
		cubitFile = fmt.Sprintf("%s/%s_cubit.dart", screenNameLower, screenNameLower)
	}

	screenData := fmt.Sprintf(`
import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

class %vScreen extends StatelessWidget {
	const %vScreen({super.key});

	@override
	Widget build(BuildContext context) {
		return BlocProvider(
		    create: (context) => %vCubit(),
		    child: BlocBuilder<%vCubit, %vVM>(
			    builder: (context, state) {
				    final cubit = context.read<%vCubit>();
				    return const Scaffold(
					    body: Center(
						    child: Text("%v")
					    )
				    );
			    }
            )
        );
	}
}
	`, screenName, screenName, screenName, screenName, screenName, screenName, screenName)

	cubitData := fmt.Sprintf(`
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:supabase_flutter/supabase_flutter.dart';

class %vCubit extends Cubit<%vVM> {
	final supabase = Supabase.instance.client;
	
	%vCubit() : super(%vVM());
	
}

class %vVM {
	%vVM.init();

	%vVM();
}
	`, screenName, screenName, screenName, screenName, screenName, screenName, screenName)

	screen := []byte(screenData)
	cubit := []byte(cubitData)

	err1 := os.WriteFile(screenFile, screen, 0644)
	check(err1)
	err2 := os.WriteFile(cubitFile, cubit, 0644)
	check(err2)
}

func createCubitFile(sceenName string) {

	cubit := fmt.Sprintf("/%s/%s_cubit.dart", sceenName, sceenName)
	f, err := os.Create(cubit)
	check(err)
	defer f.Close()

	n, err := f.WriteString(`
		Cubit.... 
	`)
	check(err)
	fmt.Printf("wrote %d bytes\n", n)
}

func flutterCLI(command string, args string) {
	cmd := exec.Command("flutter", command, args)
	switch command {
	case "create":
		b, err := cmd.CombinedOutput()
		if err != nil {
			log.Printf("Flutter function failed: %v", err)
		}
		fmt.Printf("%s\n", b)
		updatePubspec(args)
		setMainFile(args)
		setResourcesDirectory(args)
		setUI(args)
		flutterCLI("pub", "get")
	case "--version":
		b, err := cmd.CombinedOutput()
		if err != nil {
			log.Printf("Flutter function failed: %v", err)
		}
		fmt.Printf("%s\n", b)
	}
}

func updatePubspec(projectName string) {
	fileDestination := fmt.Sprintf("%s/pubspec.yaml", projectName)
	input, err := os.ReadFile(fileDestination)
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		if strings.Contains(line, "cupertino_icons: ^") {
			lines[i+1] = `
  flutter_bloc: ^8.1.3
  shared_preferences: ^2.2.2
  supabase_flutter: ^1.10.22
  email_validator: ^2.1.17
  timeago: ^3.5.0
  product_shop_ui:
    git:
      url: git@github.com:mw-productshop/fluid-ui.git
      ref: main # branch name`
		}
	}
	output := strings.Join(lines, "\n")
	err = os.WriteFile(fileDestination, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}

}

func setMainFile(projectName string) {
	fileDestination := fmt.Sprintf("%s/lib/main.dart", projectName)

	mainFileString := fmt.Sprintf(`import 'package:%s/resources/app_pages.dart';
import 'package:%s/resources/route_generator.dart';
import 'package:%s/ui/themes/dark_theme.dart';
import 'package:%s/ui/themes/light_theme.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:supabase_flutter/supabase_flutter.dart';

void main() async {
  WidgetsFlutterBinding.ensureInitialized();
  await Supabase.initialize(
    url: "",
    anonKey: ""
  );
   
  await SystemChrome.setPreferredOrientations([DeviceOrientation.portraitUp]);
  runApp(const MyApp());
}

class MyApp extends StatelessWidget {
  const MyApp({super.key});

  // This widget is the root of your application.
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Flutter Demo',
      theme: lightTheme,
      darkTheme: darkTheme,
      initialRoute: AppPages.SplashScreen,
      navigatorKey: navigationKey,
      onGenerateRoute: RouteGenerator.generateRoute,
    );
  }
}`, projectName, projectName, projectName, projectName)
	err := os.WriteFile(fileDestination, []byte(mainFileString), 0644)
	check(err)
}

func setResourcesDirectory(projectName string) {
	directory := fmt.Sprintf("%s/lib/resources", projectName)
	err := os.Mkdir(directory, 0777)
	check(err)
	createAppPages(directory)
	createNavigationRouter(directory, projectName)
}

func createAppPages(directory string) {
	appPagesData := `class AppPages {
		static const String SplashScreen = "/";
		static const String LoginScreen = "/login";
	}`
	screenFile := fmt.Sprintf("%s/app_pages.dart", directory)
	err1 := os.WriteFile(screenFile, []byte(appPagesData), os.ModePerm)
	check(err1)
}

func createNavigationRouter(directory string, projectName string) {
	routeManagerData := fmt.Sprintf(`import 'package:%s/resources/app_pages.dart';
import 'package:%s/ui/screens/auth/login/login_screen.dart';
import 'package:%s/ui/screens/auth/splash/splash_screen.dart';
import 'package:flutter/material.dart';

final GlobalKey<NavigatorState> navigationKey = GlobalKey<NavigatorState>();

class RouteGenerator {
  static Route<dynamic> generateRoute(RouteSettings settings) {
    // final args = settings.arguments;
    switch(settings.name) {
      case AppPages.SplashScreen:
        return MaterialPageRoute(builder: (_) => const SplashScreen());
      case AppPages.LoginScreen:
        return MaterialPageRoute(builder: (_) => const LoginScreen());
      default:
        return _errorRoute();
    }
  }

  static Route<dynamic> _errorRoute() {
    return MaterialPageRoute(builder: (_) {
      return Scaffold(
        appBar: AppBar(
          title: const Text('Erorr')
        ),
        body: const Center(
          child: Text('ERROR')
        ),
      );
    });
  }
}`, projectName, projectName, projectName)
	screenFile := fmt.Sprintf("%s/route_generator.dart", directory)
	err1 := os.WriteFile(screenFile, []byte(routeManagerData), os.ModePerm)
	check(err1)
}

func setUI(projectName string) {
	directory := fmt.Sprintf("%s/lib/ui", projectName)
	err := os.Mkdir(directory, 0777)
	check(err)
	setThemes(projectName)
	setInitialScreens(projectName)
}

func setThemes(projectName string) {
	directory := fmt.Sprintf("%s/lib/ui/themes", projectName)
	err := os.Mkdir(directory, 0777)
	check(err)
	createLightMode(directory)
	createDarkMode(directory)
}

func createLightMode(directory string) {
	themeData := `import 'package:flutter/material.dart';

ThemeData lightTheme = ThemeData(
	brightness: Brightness.light,
	appBarTheme: const AppBarTheme(
	backgroundColor: Colors.white,
	foregroundColor: Colors.black,
	iconTheme: IconThemeData(
		color: Colors.black
	)
	),
	bottomNavigationBarTheme: const BottomNavigationBarThemeData(
	backgroundColor: Colors.white,
	enableFeedback: false,
	type: BottomNavigationBarType.fixed,
	selectedItemColor: Colors.black,
	unselectedItemColor: Colors.grey
	),
);`
	screenFile := fmt.Sprintf("%s/light_theme.dart", directory)
	err1 := os.WriteFile(screenFile, []byte(themeData), os.ModePerm)
	check(err1)
}

func createDarkMode(directory string) {
	themeData := `import 'package:flutter/material.dart';

ThemeData darkTheme = ThemeData(
	brightness: Brightness.dark,
	appBarTheme: AppBarTheme(
	backgroundColor: Colors.grey[900],
	foregroundColor: Colors.grey[200],
	iconTheme: IconThemeData(
		color: Colors.grey[200],
	)
	),
	bottomNavigationBarTheme: const BottomNavigationBarThemeData(
	backgroundColor: Colors.black,
	enableFeedback: false,
	type: BottomNavigationBarType.fixed,
	selectedItemColor: Colors.grey,
	unselectedItemColor: Colors.white
	),
);`
	screenFile := fmt.Sprintf("%s/dark_theme.dart", directory)
	err1 := os.WriteFile(screenFile, []byte(themeData), os.ModePerm)
	check(err1)
}

func setInitialScreens(projectName string) {
	screensDirectory := fmt.Sprintf("%s/lib/ui/screens/", projectName)
	authDirectory := fmt.Sprintf("%sauth/", screensDirectory)
	loginDirectory := fmt.Sprintf("%slogin/", authDirectory)
	splashDirctory := fmt.Sprintf("%ssplash/", authDirectory)
	err := os.Mkdir(screensDirectory, 0777)
	check(err)
	err1 := os.Mkdir(authDirectory, 0777)
	check(err1)
	err2 := os.Mkdir(loginDirectory, 0777)
	check(err2)
	createScreenFile("Login", loginDirectory)
	err3 := os.Mkdir(splashDirctory, 0777)
	check(err3)
	createScreenFile("Splash", splashDirctory)
}

func check(e error) {
	if e != nil {
		if e != nil {
			log.Printf("Failed to launch failed: %v", e)
		}
		panic(e)
	}
}

func init() {
	FlutterCmd.AddCommand(versionCmd)
	FlutterCmd.AddCommand(createCmd)
	FlutterCmd.AddCommand(createScreen)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// flutterCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// flutterCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
