'use strict';

app.config(function ($stateProvider, $urlRouterProvider) {

  // ROUTING with ui.router
  $urlRouterProvider.otherwise('/');
  $stateProvider
    // this state is placed in the <ion-nav-view> in the index.html
    .state('main', {
      url: '/',
      abstract: true,
      // controller: 'HackCtrl',
      resolve: {
        allSnippets: function (Snippets) {
          return Snippets.getUberAll();
        },
        allBuckets: function (Buckets) {
          return Buckets.fetchAll();
        },
        resolvedFS: function (FS) {
          return FS.fetchFS();
        }
      }
    })
    .state('main.hackon', {
      url: '/:snippetId',
      
    });
});

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