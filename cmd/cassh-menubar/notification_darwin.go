//go:build darwin

package main

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation -framework UserNotifications

#import <Foundation/Foundation.h>
#import <UserNotifications/UserNotifications.h>

void requestNotificationPermission() {
    UNUserNotificationCenter *center = [UNUserNotificationCenter currentNotificationCenter];
    [center requestAuthorizationWithOptions:(UNAuthorizationOptionAlert | UNAuthorizationOptionSound)
                          completionHandler:^(BOOL granted, NSError * _Nullable error) {
        if (error) {
            NSLog(@"Notification permission error: %@", error);
        }
    }];
}

void sendNativeNotification(const char *title, const char *body) {
    NSString *nsTitle = [NSString stringWithUTF8String:title];
    NSString *nsBody = [NSString stringWithUTF8String:body];

    UNUserNotificationCenter *center = [UNUserNotificationCenter currentNotificationCenter];

    UNMutableNotificationContent *content = [[UNMutableNotificationContent alloc] init];
    content.title = nsTitle;
    content.body = nsBody;
    content.sound = [UNNotificationSound defaultSound];

    NSString *identifier = [[NSUUID UUID] UUIDString];
    UNNotificationRequest *request = [UNNotificationRequest requestWithIdentifier:identifier
                                                                          content:content
                                                                          trigger:nil];

    [center addNotificationRequest:request withCompletionHandler:^(NSError * _Nullable error) {
        if (error) {
            NSLog(@"Failed to send notification: %@", error);
        }
    }];
}
*/
import "C"
import "unsafe"

func initNotifications() {
	C.requestNotificationPermission()
}

func sendNativeNotification(title, body string) {
	cTitle := C.CString(title)
	cBody := C.CString(body)
	defer C.free(unsafe.Pointer(cTitle))
	defer C.free(unsafe.Pointer(cBody))

	C.sendNativeNotification(cTitle, cBody)
}
