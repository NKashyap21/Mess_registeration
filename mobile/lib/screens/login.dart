import 'package:flutter/material.dart';
import 'package:flutter_frontend/screens/home.dart';
import 'package:flutter_frontend/widgets/button.dart';
import 'package:flutter_frontend/theme/app_theme.dart';
import 'package:flutter_frontend/api/google_sign_in_service.dart';
import 'package:flutter_frontend/api/api_service.dart';

class LoginScreen extends StatefulWidget {
  const LoginScreen({super.key});

  @override
  State<LoginScreen> createState() => _LoginScreenState();
}

class _LoginScreenState extends State<LoginScreen> {
  bool _loading = true;
  bool _signingIn = false;

  @override
  void initState() {
    super.initState();
    _checkAutoLogin();
  }

  Future<void> _checkAutoLogin() async {
    final token = await ApiService.getToken();
    if (token != null && token.isNotEmpty) {
      Navigator.pushReplacement(
        context,
        MaterialPageRoute(builder: (context) => const HomeScreen()),
      );
    } else {
      setState(() {
        _loading = false;
      });
    }
  }

  void _handleSignIn() {
    if (_signingIn) return;

    setState(() {
      _signingIn = true;
    });

    signInWithGoogle()
        .then((user) {
          if (user) {
            if (mounted) {
              Navigator.pushReplacement(
                context,
                MaterialPageRoute(builder: (context) => const HomeScreen()),
              );
            }
          } else {
            if (mounted) {
              ScaffoldMessenger.of(context).showSnackBar(
                const SnackBar(
                  content: Text('Google Sign-In failed. Please try again.'),
                  backgroundColor: Colors.red,
                ),
              );
            }
          }
        })
        .catchError((e) {
          if (mounted) {
            ScaffoldMessenger.of(context).showSnackBar(
              SnackBar(
                content: Text('Error: ${e.toString()}'),
                backgroundColor: Colors.red,
              ),
            );
          }
        })
        .whenComplete(() {
          if (mounted) {
            setState(() {
              _signingIn = false;
            });
          }
        });
  }

  @override
  Widget build(BuildContext context) {
    if (_loading) {
      return const Scaffold(body: Center(child: CircularProgressIndicator()));
    }

    return Scaffold(
      body: SafeArea(
        child: LayoutBuilder(
          builder: (context, constraints) {
            return SingleChildScrollView(
              child: ConstrainedBox(
                constraints: BoxConstraints(minHeight: constraints.maxHeight),
                child: IntrinsicHeight(
                  child: Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 16.0),
                    child: Column(
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: [
                        const Spacer(),
                        _buildLoginContent(context),
                        const Spacer(),
                        _buildFooter(),
                      ],
                    ),
                  ),
                ),
              ),
            );
          },
        ),
      ),
    );
  }

  Widget _buildLoginContent(BuildContext context) {
    return Column(
      mainAxisAlignment: MainAxisAlignment.center,
      children: [
        Image.asset(
          'assets/images/mess_icon.png',
          width: 50,
          height: 50,
          fit: BoxFit.contain,
        ),
        const SizedBox(height: 16),
        const Text('Welcome to', style: TextStyle(color: AppColors.lightGrey)),
        const SizedBox(height: 8),
        const Text(
          'IITH Mess Portal',
          style: TextStyle(
            fontSize: 24,
            fontWeight: FontWeight.bold,
            color: AppColors.white,
          ),
          textAlign: TextAlign.center,
        ),
        const SizedBox(height: 32),
        AppButton.icon(
          label: _signingIn ? 'Signing in...' : 'Sign in with Google',
          onPressed: _handleSignIn,
          icon: _signingIn
              ? const SizedBox(
                  width: 24,
                  height: 24,
                  child: CircularProgressIndicator(
                    strokeWidth: 2,
                    valueColor: AlwaysStoppedAnimation<Color>(Colors.white),
                  ),
                )
              : Image.asset(
                  'assets/images/google_icon.png',
                  width: 24,
                  height: 24,
                  fit: BoxFit.contain,
                ),
          backgroundColor: AppColors.darkGrey,
          borderRadius: BorderRadius.circular(0),
        ),
      ],
    );
  }

  Widget _buildFooter() {
    return Container(
      padding: const EdgeInsets.symmetric(vertical: 16.0),
      child: const Row(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Text(
            'Brought to you with ',
            style: TextStyle(color: AppColors.lightGrey),
          ),
          Icon(Icons.favorite, color: Colors.deepPurple, size: 16),
          Text(' by Lambda', style: TextStyle(color: AppColors.lightGrey)),
        ],
      ),
    );
  }
}
