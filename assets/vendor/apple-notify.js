// https://developer.apple.com/library/mac/documentation/NetworkingInternet/Conceptual/NotificationProgrammingGuideForWebsites/LocalNotifications/LocalNotifications.html#//apple_ref/doc/uid/TP40012932-SW1
var notify = function (notificationObject) {
    // Check for notification compatibility.
    if (!'Notification' in window) {
        // If the browser version is unsupported, remain silent.
        return;
    }
    // Log current permission level
    console.log(Notification.permission);
    // If the user has not been asked to grant or deny notifications
    // from this domain...
    if (Notification.permission === 'default') {
        Notification.requestPermission(function () {
            // ...callback this function once a permission level has been set.
            notify();
        });
    }
    // If the user has granted permission for this domain to send notifications...
    else if (Notification.permission === 'granted') {
        var n = new Notification(
                    notificationObject.text,
                    {
                      'body': notificationObject.body,
                      // ...prevent duplicate notifications
                      'tag' : notificationObject.tag
                    }
                );
        // Remove the notification from Notification Center when clicked.
        n.onclick = function () {
            this.close();
        };
        // Callback function when the notification is closed.
        n.onclose = function () {
            console.log('Notification closed');
        };
    }
    // If the user does not want notifications to come from this domain...
    else if (Notification.permission === 'denied') {
        // ...remain silent.
        return;
    }
};

// Original
// ------
// var notify = function () {
//     // Check for notification compatibility.
//     if (!'Notification' in window) {
//         // If the browser version is unsupported, remain silent.
//         return;
//     }
//     // Log current permission level
//     console.log(Notification.permission);
//     // If the user has not been asked to grant or deny notifications
//     // from this domain...
//     if (Notification.permission === 'default') {
//         Notification.requestPermission(function () {
//             // ...callback this function once a permission level has been set.
//             notify();
//         });
//     }
//     // If the user has granted permission for this domain to send notifications...
//     else if (Notification.permission === 'granted') {
//         console.log('notifying?!');
//         var n = new Notification(
//                     'New message from Liz',
//                     {
//                       'body': 'Liz: "Hi there!"',
//                       // ...prevent duplicate notifications
//                       'tag' : 'unique string'
//                     }
//                 );
//         // Remove the notification from Notification Center when clicked.
//         n.onclick = function () {
//             this.close();
//         };
//         // Callback function when the notification is closed.
//         n.onclose = function () {
//             console.log('Notification closed');
//         };
//     }
//     // If the user does not want notifications to come from this domain...
//     else if (Notification.permission === 'denied') {
//         // ...remain silent.
//         return;
//     }
// };

