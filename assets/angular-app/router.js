'use strict';

app.config(function ($stateProvider, $urlRouterProvider, $urlMatcherFactoryProvider, $locationProvider) {
  // $urlMatcherFactoryProvider.strictMode(false);
  // $locationProvider.html5Mode(true);
  $urlRouterProvider.otherwise('/');
  $stateProvider
    .state('main', {
      url: '/',
      resolve: {
        stuff: function (Snippets, Buckets, FS, Utils, $state, $stateParams) {
          Snippets.getUberAll()
            .then(function (res) {
              Snippets.setManyToSnippetsLib(res.data);
            })
            .then(function () {
              return Buckets.fetchAll();
            })
            .then(function (res) {
              Buckets.storeManyBuckets(res.data);
            })
            .then(function () {
              return FS.fetchFS();
            })
            .then(function (res) {
              FS.storeFS(res.data);
            })
            .then(function () {
              if (Utils.typeOf(Snippets.getSnippetsLib()) !== 'null') {
                var mostRecent = Snippets.getMostRecent(Snippets.getSnippetsLib());
                console.log('going to state: ' + mostRecent.id);
                $state.go('main.hackon', {snippetId: mostRecent.id});
              } else {
                var defaulty = Config.DEFAULTSNIPPET;
                Snippets.setOneToSnippetsLib(defaulty);
                console.log('going to state: ' + defaulty.id);
                $state.go('main.hackon', {snippetId: defaulty.id});
              }
            })
        }
      }
    })
    .state('main.hackon', {
      url: ':snippetId',
      controller: 'HackCtrl',
      resolve: {
        bucket: function (Snippets, Buckets, $stateParams, $log) {
          console.log('resolving bucket');
          if (Object.keys(Snippets.getSnippetsLib()).length === 0 && Snippets.getSnippetsLib().constructor === Object) {
            var buck = Buckets.getBuckets()[0]; // presuming bucket is snippets?
            // $scope.data.cs.bucketId = $scope.data.cb.id; // set default snippet to have default bucket id
          } else {
            var buck = Buckets.getBuckets()[Snippets.getSnippetsLib()[$stateParams.snippetId].bucketId]; // set current bucket to be current snippet's bucket
          }
          $log.log('buck -> ', buck);
          return buck;
        },
        snippet: function (Snippets, $stateParams, $log) {
          var snip = Snippets.getSnippetsLib()[$stateParams.snippetId];
          $log.log('snip -> ', snip);
          return snip;
        }
      }
    });
});

// 'use strict';

// app.config(function ($stateProvider, $urlRouterProvider) {

//   // ROUTING with ui.router
//   $urlRouterProvider.otherwise('hack/');
//   $stateProvider
//     // this state is placed in the <ion-nav-view> in the index.html
//     .state('main', {
//       url: '/hack/',
//       abstract: true,
//       // controller: 'HackCtrl',
//       resolve: {
//         // TODO handle ajax errors.
//         allSnippets: function (Snippets, Utils, $state, $stateParams) {
//           console.log("i am routerring!");
//           return Snippets.getUberAll().then(function (res) {
//             if (Utils.typeOf(res.data) !== 'null') {
//               Snippets.setManyToSnippetsLib(res.data);
//               var mostRecent = Snippets.getMostRecent(Snippets.getSnippetsLib());
//               $state.go(main.hackon({snippetId: mostRecent.id}));
//             } else {
//               var defaulty = Config.DEFAULTSNIPPET;
//               Snippets.setOneToSnippetsLib(Config.DEFAULTSNIPPET);
//               $state.go(main.hackon({snippetId: defaulty.id}));
//             }
//           });
//         },
//         allBuckets: function (Buckets, Snippets) {
//           return Buckets.fetchAll().then(function (res) {
//             Buckets.storeManyBuckets(res.data);
//           });
//         },
//         resolvedFS: function (FS) {
//           return FS.fetchFS().then(function (res) {
//             FS.storeFS(res.data);
//           });
//         }
//       }
//     })
//     .state('main.hackon', {
//       url: '/:snippetId',
//       controller: 'HackCtrl',
//       resolve: {
//         bucket: function (Snippets, Buckets, $stateParams) {
//           if (Object.keys(Snippets.getSnippetsLib()).length === 0 && Snippets.getSnippetsLib().constructor === Object) {
//             return Buckets.getBuckets()[0]; // presuming bucket is snippets?
//             // $scope.data.cs.bucketId = $scope.data.cb.id; // set default snippet to have default bucket id
//           } else {
//             return Buckets.getBuckets()[Snippets.getSnippetsLib()[$stateParams.snippetId].bucketId]; // set current bucket to be current snippet's bucket
//           }
//         },
//         snippet: function (Snippets, $stateParams) {
//           return Snippets.getSnippetsLib()[$stateParams.snippetId];
//         }
//       }
//     });
// });

//
//      })
      // .state('main.listDetail', {
      //   url: '/list/:poopId',
      //   views: {
      //     'main-list': {
      //       templateUrl: 'main/templates/list-detail.html',
      //       controller: 'ListDetailCtrl',
      //       resolve: {
      //         poop: function ($stateParams, Tallyrally) {
      //           return Tallyrally.get({id: $stateParams.poopId}, function (res) {
      //             return res;
      //           });
      //         }
      //         // poopId: function ($stateParams) {
      //         //   return $stateParams.poopId;
      //         // }
      //       }
      //     }
      //   }
      // })