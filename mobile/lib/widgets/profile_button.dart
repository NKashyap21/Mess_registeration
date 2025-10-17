import 'package:flutter/material.dart';
import 'package:flutter_frontend/theme/app_theme.dart';
import 'package:flutter_frontend/screens/login.dart';
import 'package:flutter_frontend/api/google_sign_in_service.dart';

class ProfileButton extends StatelessWidget {
  final Map<String, dynamic>? userInfo;

  const ProfileButton({Key? key, this.userInfo}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      onTap: () => _showProfileModal(context),
      child: Padding(
        padding: const EdgeInsets.only(right: 16.0),
        child: CircleAvatar(
          backgroundColor: AppColors.darkGrey,
          radius: 20,
          child: const Icon(Icons.person, color: AppColors.white),
        ),
      ),
    );
  }

  void _showProfileModal(BuildContext context) {
    showDialog(
      context: context,
      builder: (BuildContext context) {
        return Dialog(
          backgroundColor: Colors.black,
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(8.0),
          ),
          child: Padding(
            padding: const EdgeInsets.all(16.0),
            child: Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                Row(
                  mainAxisAlignment: MainAxisAlignment.end,
                  children: [
                    IconButton(
                      icon: const Icon(Icons.close, color: AppColors.white),
                      onPressed: () => Navigator.of(context).pop(),
                      padding: EdgeInsets.zero,
                      constraints: const BoxConstraints(),
                    ),
                  ],
                ),
                const SizedBox(height: 8),

                // Profile info
                const CircleAvatar(
                  backgroundColor: AppColors.darkGrey,
                  radius: 32,
                  child: Icon(Icons.person, color: AppColors.white, size: 32),
                ),
                const SizedBox(height: 12),

                if (userInfo != null) ...[
                  Text(
                    userInfo!['name']?.toString() ?? 'Unknown User',
                    style: const TextStyle(
                      color: AppColors.white,
                      fontSize: 16,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                  const SizedBox(height: 4),
                  Text(
                    userInfo!['roll_number']?.toString() ?? 'No roll number',
                    style: const TextStyle(
                      color: AppColors.lightGrey,
                      fontSize: 14,
                    ),
                  ),
                  if (userInfo!['id'] != null) ...[
                    const SizedBox(height: 4),
                    Text(
                      'ID: ${userInfo!['id']}',
                      style: const TextStyle(
                        color: AppColors.lightGrey,
                        fontSize: 12,
                      ),
                    ),
                  ],
                ] else ...[
                  const Text(
                    'Failed to load user info',
                    style: TextStyle(color: AppColors.orange, fontSize: 14),
                  ),
                ],
                const SizedBox(height: 20),

                InkWell(
                  onTap: () async {
                    final navigator = Navigator.of(context);
                    final currentContext = context;

                    navigator.pop();

                    showDialog(
                      context: currentContext,
                      barrierDismissible: false,
                      builder: (dialogContext) =>
                          const Center(child: CircularProgressIndicator()),
                    );

                    try {
                      await signOut().timeout(
                        const Duration(seconds: 10),
                        onTimeout: () {},
                      );
                    } catch (e) {
                    } finally {
                      try {
                        navigator.pop();
                      } catch (e) {}

                      try {
                        navigator.pushAndRemoveUntil(
                          MaterialPageRoute(
                            builder: (context) => const LoginScreen(),
                          ),
                          (route) => false,
                        );
                      } catch (e) {
                        print('Error navigating to login: $e');
                      }
                    }
                  },
                  child: Container(
                    padding: const EdgeInsets.symmetric(
                      vertical: 12,
                      horizontal: 16,
                    ),
                    width: 120,
                    decoration: BoxDecoration(
                      color: AppColors.darkGrey,
                      borderRadius: BorderRadius.circular(8),
                    ),
                    child: const Center(
                      child: Text(
                        'Sign Out',
                        style: TextStyle(color: AppColors.white, fontSize: 14),
                      ),
                    ),
                  ),
                ),
              ],
            ),
          ),
        );
      },
    );
  }
}
