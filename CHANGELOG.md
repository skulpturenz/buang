# Changelog

All notable changes to this project will be documented in this file.

## 1.0.0-canary.3 - 2026-02-04

[b1ae734](b1ae734e5c4c0435ccb39a886cb31cc5e319da9d)...[316b4c9](316b4c9f6ba41c5b2a4d57e2d45773545bf1ed96)

### Miscellaneous Tasks

- Update release workflow (#14) ([200ff8d](200ff8d734e27656a6ae18a63ced5129a6ea3a80))

## 1.0.0-canary.2 - 2026-02-04

[df15f85](df15f85203a298344a21dcf54a10261ff99dc326)...[b1ae734](b1ae734e5c4c0435ccb39a886cb31cc5e319da9d)

### Miscellaneous Tasks

- Cleanup ([c685c31](c685c31dd287ed4904d3c09bb2dbd7b5b4b3179c))
- Add examples ([7e5918e](7e5918e2e3c9ee8d55c5e2835639fc161e2e66a4))
- Add status checks workflow ([b0dfe9a](b0dfe9a8a65191a3bbc5e3f2746f74e6839d522e))
- Test workflow secrets (#12) ([6ceba7d](6ceba7daa3acdaa14d662c5a08e9b8e36dacb117))
- Update changelog ([b1ae734](b1ae734e5c4c0435ccb39a886cb31cc5e319da9d))

## 1.0.0-canary.1 - 2026-02-03

[33acf8d](33acf8d51259a0890af3500e45f662747906cfe1)...[df15f85](df15f85203a298344a21dcf54a10261ff99dc326)

### Miscellaneous Tasks

- Update changelog format, update release workflow ([4c941ab](4c941ab0a5a17d6e2b20c885b341890be7753586))
- Update changelog ([df15f85](df15f85203a298344a21dcf54a10261ff99dc326))

## 1.0.0-canary.0 - 2026-02-03

### Bug Fixes

- Issues with sqlite, issues with SelectProjectsDesc ([d251f1d](d251f1d3f243e466b3048b573c74ef5c34bf4275))
- Add req validations, buang all active deployments by branch, schema updates, component updates ([9f2571f](9f2571f3ed924b534e7aff94f8cdd4546156eb8c))
- Issues with deployment and validation ([52e7d45](52e7d45ccf057e62adcd3851032a186b63cdfc7b))
- Update ServiceEntrypoint validation ([af0532c](af0532ce442a57f48a0d16043a1c5dab0b1c8c88))
- Setting env vars ([b043204](b043204ac28e3b69984076b42f5ff65f0719636d))
- Add buang_branch, fix deployment issues ([0ee43f2](0ee43f20619425f62d12dd4eca582f8c5daf8ba0))
- Update clone configuration ([d61c27d](d61c27d5fd4463c932fe0e3d5c6fa383029419e3))
- Pg SelectProjectsDesc ([4c14716](4c14716316b2a2a27aa95a8d48a95e41b950b682))
- Run housekeeping worker ([9dfd059](9dfd05949a42d4c17bc4ddb960529634e6640af0))
- Dbos issues ([1d8cd24](1d8cd24cf4b53a2dbcc2611fdf47cc718a0a9344))
- Housekeeping workflow DI ([1da26cc](1da26cc42c028b901f0d9b46c5b0083b75815e86))
- Panic recovery and improve ([958e3bd](958e3bd21d920521e4c94f342a71d03a257cab7a))
- Handle compose down active deployments errors better. spamming the execute button 60+ times surfaces it. services are still running but compose file is gone: ostrich, find really old services and kill them or maybe not using temp dirs is more suitable ([50c0b83](50c0b83e2467bb03f04ac71baa84edc978d35343))
- Housekeeping workflow nil ptr, concurrent deployments breaking 1 active deployment constraint ([a75783e](a75783e556f961f370fcdc5d3fe2801b10b95ac4))
- Status code ([8c930ec](8c930ec857d719ac98b80d42fd843cc6b913ba79))
- 1 active deployment per branch for every project ([5db9801](5db9801383540ca698cc35eeebb1a123c0b16064))
- Housekeeping temporal workflow ([4cab103](4cab10364d51eb809af818b37d3113d28553687b))
- Stringify stack before saving to db, fix issue with sqlite CreateDeployment ([695253c](695253c6af1ca52c0d64ecd0e9540b193be77c10))
- Buang_deployment activity panic if clonePath is nil ([33ff451](33ff4513750dd70f51945a1361f220d1fac9a89b))
- Nil envs ([e478dd9](e478dd91aab5a94a8bf1e0b16d7a6665a32b9234))
- Handle nil map cases at db wrappers ([a18b800](a18b80088de9d8d4cd83bd893672fc0911f32607))
- DeleteProject, UpdateProject swagger annotations, update create_deployment query to check if project has been deleted ([37e8928](37e8928f40748ffed3cf546c50d42147636d5e65))
- GetDeploymentLogs returns too early when deployment is already in progress but not completed ([8207dd1](8207dd114912fec1f05546e3a1794d406a9253c4))
- Update isDeploying check ([e71860c](e71860cc3faf1ff0ebfebcdaf3cd76d517b5a273))
- Deployment router paths ([becdba4](becdba448fd6d50d4ece6cf8de6220bd144463de))
- On delete cascade ([529ac1b](529ac1b883212b3f22cb4661a75228b4e2aa6e79))
- Update cpu usage calculation ([8198c98](8198c9829f9c611454df9259d6002ad6c7efd84e))
- MemUsage calculation, used total_inactive_file vs inactive_file ([7386738](738673888f666ad7578bd52b70a49ffb4a277dfa))
- Sqlite in memory getting cleared ([a2eb58d](a2eb58d72fe93327f51cff9a4d4319aa5243d83e))
- Compose logs not included with temporal ([3361f79](3361f795330c4b9b2960e3bc38c394bba11dec78))
- Retries and handle context timeouts ([5c6cec0](5c6cec01d9d68528922bd89201917ec8475a7b69))
- Increase step retries ([2b91990](2b91990dbeab12daf658ae2eeee9416a76cec474))

### Documentation

- Swagger annotations ([65c54f9](65c54f91fdf69aae4fc4bed8590895970763de35))
- Update CreateDeployment swagger ([7d8c47e](7d8c47e74a0d427627229e103a2d0487101e84b5))
- Add TODO ([e16fd23](e16fd2338be9a6f61e9c274c208394cc9ee0c15b))

### Features

- Add queries to create, select and update deployments ([a4e3170](a4e317090013243d33602f7f1236d2e764e5e406))
- Add deployments component ([b2a8697](b2a8697ffdf2a9fc23c6b26b053428449a9665f2))
- Add routes for deployments ([bbd4101](bbd4101285b567e887d7af384792e9d869e04592))
- Add docker, fs and git components ([3882285](3882285c437122e6311d72672d69406d969326db))
- Add workers, activities, update schema ([495f8ec](495f8ec8047748251f8cb86638754726ca25a483))
- Add git authn, update taskfile, other small fixes ([d8e4902](d8e4902b20a860539fc5370bd519cad22f260c5b))
- Buang all active deployments when a new one occurs ([05e1dae](05e1daeb8867249f0f6d486b8b912d58adff2da3))
- List deployments and projects req validation ([bfc8252](bfc825208c576e724ca8f82afd13ca10b99f723a))
- Housekeeping worker ([7cbac25](7cbac25e5058405a6bc9f11e004e1854347b9050))
- Add buang branch endpoint, register housekeeping worker ([d64e523](d64e523f2ab7dd793b9cb39f63c7e41e8dd68fb2))
- Add endpoint to find project by repository, enforce only one active repository at a time, add github actions examples (untested) ([138a3bf](138a3bf0f3b740a6bcacb39dd200af5b591c8b3b))
- Add FindDeploymentById endpoint, update actions to leave a link to the deployed environment ([15d05e2](15d05e2077d80432bbbfb533c3e3eaa7817eab5b))
- Docker prune every 2 days ([7544e70](7544e70e1b3157a4771236409c9f1754b5633c8e))
- Add support for dbos ([f9d2a46](f9d2a46b4f2567aaf6573f5a558d47a81740c7ad))
- Recover from dbos workflow failures ([77a9e05](77a9e053450c4987969318a255cb36f126a1cec0))
- Update deployment status to error after recovering from panic ([adce192](adce1921848b9714be6e3bc55c002d9fa8a40633))
- More panic handling ([6720e1d](6720e1dcc606801fc81da7ba54af401b732334cc))
- Allow connecting to dbos conductor ([10ddf05](10ddf057720955584bb4c696d44c758541b3edfb))
- Deployment logs ([d3ff88b](d3ff88b27bc38c3b65213190bad3e637de8c66a9))
- Allow blocking until deployment completes ([9cd2d5f](9cd2d5f774c92232d1acb4bcb228b32522261d9d))
- Get deployment logs endpoint ([62eac91](62eac9123b28950100eff17f04c085e050f479bf))
- Allow waiting for deployment to complete when retrieving logs ([52cd643](52cd64391f7b5dd45ece8bacddc0c06a4a26cb52))
- Git clone logs ([91905fe](91905fe575183a156b86c26b3f5f354b9383dd83))
- Provide BUANG_DEPLOYMENT_PATH when composing up ([83860c7](83860c7c0f8958e3c172933056b05f700c4a9005))
- Add queries DeleteProject and UpdateProject ([df565bf](df565bfbdf91a35697d4009e5315647a539517bf))
- Add diagnostic log table ([a73ea76](a73ea760d4aa27240babff573b4fb545e4e8c529))
- Bootstrap ([a8844fa](a8844fa36324bb716e11b0fb80ce05bb1f0ec39d))
- Capture all panics from workflows ([fa582f6](fa582f6c5debe4a79a3e78479c61b6eae7e26d5f))
- Prune diagnostic logs ([be24caa](be24caab5be1111ac4da96b82d0e0eb1043e9df5))
- Parse stack ([b0d1250](b0d12506ab83cb421a2f9c6774cfda6a0d75c22c))
- Enable sticky cookies ([b2e4951](b2e495169ce2ea6721265a4d7323925252119ccf))
- DeleteProject, UpdateProject ([e1d90a2](e1d90a25f8d47889bf4cc2b4c8d466fcbe9e5fee))
- Stream deployment logs ([302f984](302f984216d552fc649fd52da8e41a6ed8d4588e))
- Fallback to waiting for deployment to complete if streaming is unsupported ([263d1b4](263d1b4b65e96f45641372c938e910ebbe8f16ca))
- Docker stats ([07a6123](07a6123848fa56284994a73db719f1c4e21f0c50))
- Sort docker stats ([aae360e](aae360eb0a491e96470d33f8f575ef493f6a0fb8))
- Add version route ([517ee8c](517ee8ca50724f915c7e218b89d0120facc7196e))
- Add image and image id to stats response ([28db4bf](28db4bf113803ce919a9ec6bd842a6077ed80a29))
- Enable assertions in dev, ensure all singletons are initialized ([c25c679](c25c67936a66740661c5c7e955a08dd64c651328))
- Return container status and state for stats, GET /diagnostics/version -> GET /diagnostics/info - add tz and uptime information ([27586b9](27586b9bae96b9be9ccdc7d9b4baedf41aed60f0))
- Update compose up build options to push, pass writer for build progress ([7a0dcde](7a0dcde9fc304f1ada4f3cc6d9d0ce2c02475e3d))
- Add fields for containers in DurableExecutorConfiguration ([f545c09](f545c0936928c0754d7bf54cec7d99b61b82f9b2))

### Miscellaneous Tasks

- Initial commit ([b4f0551](b4f05518e0cd03ac742f756a5059db59bdf4a10b))
- Setup project ([2526f95](2526f95b69e56d2d8d7ce306919e4e397ca1e3ca))
- Db models ([fd7fe1a](fd7fe1aa3cd0b4f09df2c1649ab8696e9bf9931a))
- Setup db connections, di, slog, refactor ([66333cd](66333cdbd130e0834c36f8b15292faea44d28758))
- Update schema ([dce5a4c](dce5a4c97b5af15ad946a403c7b95e34ceb9834c))
- Add Taskfile ([4392b7f](4392b7f2c1bfa91888b90ae26553c8c76024dc78))
- Update swagger endpoint security, update bundled migrations ([1634ddc](1634ddcc872614cd3e53b7dbe46940a40340735a))
- Task generate-all ([305d5ef](305d5efb522e59f1888bb328eff7104b43ddaa72))
- Add infrastructure ([1cdff82](1cdff8230d4bb8e8bc03766a9ace2cdcfe77b1d8))
- Update README ([b427e87](b427e870d0b074addef7e8e1bac2c50a9a0fbf3f))
- Add demo to README ([3793a64](3793a64546a76cb50700677ee28b4c7697fd029c)), Signed-off-by:Naveen Mathew <55116576+nmathew98@users.noreply.github.com>
- Update README ([747828c](747828c7ea2997a2377a23b668e88afcba55f192))
- Add deploy-prod task and prod compose file ([c1b9757](c1b9757573439d69531bf823cce388ffebc0b22e))
- Update README ([1d647f3](1d647f322b741486cd38f3545f73bd215acb7952))
- Task generate-all ([4895076](4895076e3c0bff898f46023b64a482e01725b2c7))
- Update README ([f94bb45](f94bb453d113aff6b647ce34bd8a367d4b371f1a))
- Add configuration.md docs ([272c0f9](272c0f934560d20f29ab9518a0270d83836629e2))
- Add git-chglog ([72dc3cc](72dc3cc0d7de5b40441fba1360b96f0f7439c9e4))
- Update docker compose to add health check and update volume mount points, still untested ([471d6ab](471d6ab26469961896b9c02e98b28580e9485b26))
- Panic if BUANG_VERSION is not set in production ([22f8e7a](22f8e7a651500fe7e8a3c8e96055d0e19139ada0))
- Docker in docker for buang ([a4d0fe3](a4d0fe37e0d9d796af01cf99273feeee57f5f5b7))
- Update Dockerfile ([3e0c0d9](3e0c0d9ba87af8118c656874b7c01531dbff84e0))
- Setup utils for integration tests ([b82de86](b82de868319ecaf7e08c3c8019fc2eb772143487))
- Update taskfile ([96a1d8a](96a1d8ae85872ae9b9f7f1d79cdb499a39117ef6))
- Update logger configuration ([f2f8887](f2f8887fa2e764af4865206e1e892dde4a68a810))
- Fix Dockerfile ([985a489](985a48952e125aa12c7a30dc7fa1eeab461e7054))
- Add test github action workflow ([1ecf759](1ecf75959a140b767d93ba4121ddb93593b4f7c9))
- Update gh actions workflow to cache go package downloads by default ([fda77ca](fda77ca0064a302b32fe151e091f95cf11407703))
- Fix cache dependency path ([a2f722c](a2f722cf35c589b279a71b76572c1dd71efb6601))
- Authn github test-services action ([9eb19e0](9eb19e01a217ae7be3634872ff25990c6c9d3161))
- Authenticate install task step in test services gh workflow ([324f7da](324f7daac9c9630a4df5e0ba4db54b4601fa980a))
- Migrate to git-cliff ([9b25972](9b25972bef6c5ea34b40a22ae551300735edd7fe))
- Add release workflow ([19faf3c](19faf3c51ed5b119407edf50a6b2665059512fec))
- Update changelog ([33acf8d](33acf8d51259a0890af3500e45f662747906cfe1))

### Refactor

- Update deployment routes, int temporal and http ([cf1c8fc](cf1c8fc90e571c71d9a878153b2a3d42b816e657))
- Assertions, update Taskfile ([03288a1](03288a126ca898f9f35c4a652fad0263a8742913))
- Don't treat panics as errors, the application is setup wrong if it panics. iirc if compose tries to open a directory and it does not exist then we have a panic, but i think we shouldn't consider it as a deployment failure because the data is right, the issue is that we're using /tmp for files which are temporary but expected to be there because we don't try to retrieve the files again if its not once we have initially done it ([6b363a3](6b363a3c9836220c5f93786cebfb23645b58827c))
- More err handling ([67459de](67459dea340b56035eab23ba4f8f534ed5e86516))
- Sqlc named params for pg, update queries UpdateDeployment and UpdateDeploymentStatus to check for project id ([45bcba0](45bcba043318316404ed8664f4585e864cfe8936))
- Rename sqlite constraints ([79f0052](79f0052cead2a9b52d8e1bb131ec0ea578dc6aa7))
- Envs ([e14aaec](e14aaecc83aea043262e09cf06b6630029b86d99))
- Standardize diagnostic logs a little ([b11b162](b11b1624f3a8988716aebe6cb7a4598223acb628))
- Dberrors ([0040fb0](0040fb0ff16f4c8c2428755fdee49202fa82f9eb))
- Tidy GetDeploymentLog and PollDeploymentLog ([a590b13](a590b13a900367095d61cb897e50b1cc977983c3))
- Simplify and standardize workflow panic recovery ([09c779f](09c779f416d8e7488bc08b9fe4b9bcda8f4a78ad))
- Reorganize handlers ([7c4dad2](7c4dad2379330d822de7569bd3ef314c888f42ba))
- Update swagger grouping ([5ff13cb](5ff13cb1575e7012ade81a722375412061ecf867))
- Collect ([2bd5471](2bd5471eb86a677d0ffe07929cba09ee9b0dbd69))
- Update DI - inject docker client, inject gorilla schema from main ([32371cb](32371cb77fa1c23fab0bea0d61ac8f50f0cd0324))
- Scope envs, allow configuring traefik ports ([e02e81a](e02e81a5470ea1ce63df41bd872be219e62c6979))
- Make durable executor fields in ApplicationServices private ([a37878e](a37878e541ccfd7ae085b664979f0805e51fef62))
- Replace conditionals with polymorphism ([2aac217](2aac217565809013dbc656a7a5932335e5ce5760))
- Ensure clients implement the Workflows interface ([2f8d7e0](2f8d7e06a52830779b1f9e78b82123e33e370534))
- Improve sagas for concurrency ([94f3fc3](94f3fc32539bffcf23369346872cf7f21f54764b))
- Improve DI for services which don't have global scope ([cd06cd3](cd06cd385d7ce5af5b1d99f423146606430dee58))

### Testing

- New project with deployment ([32685b4](32685b4c317c84cdcda70f179ed298c95afd34f3))
- Check sha and branch of checkout ([68cc4fb](68cc4fbd0d5a31bd3a0f08ee24761cf1be4e7f11))
- Check compose service after deployment ([eb96ee2](eb96ee25fc73d871681d37e6a0eb04f590e8eb37))
- Check envs are set correctly ([d09c482](d09c482653ce6c7db797d6177d1e7aab65630926))
- Listing projects, deployments, paging, buang branch ([fdf23bf](fdf23bf86c218e3f62395d2d2e1dabcf69c1e3ba))
- Docker stats. and fix docker stats sorting, we were sorting it in ascending and when we returned the response we were iterating from newest so we end up with a response which was in descending. also update sorting tests in list_deployments_rest and list_projects_test ([fcee3ef](fcee3efb70ec654a9bc3fa402bb0783a2ef1b0f5))
- Supress testcontainer logs ([378f357](378f3570f7c5716c4fb4afb261394d9081c71c00))
- Run tests in parallel, 45s max, under 30 most of the time. slightly flaky but more when runnign in watch mode. there are orphan containers ([45298be](45298be577fdee3d77b160338093ca87507b550f))
- Retry failing tests ([167ae82](167ae82539508048c5c154eae8593f515faf5ede))
- Add TemporalPg test config ([884298e](884298e4446cd40150ec9276ac620812f3dbcabc))
- Assert busiest and idlest containers ([53dde98](53dde98a5536742363b0725b817368e99820d45d))
- Remove time.Sleep, better logging for retry ([20bff40](20bff404b95ff4e9ebb33545188c7cfe30b255b6))
- TemporalPg for TestListProjectsExcludesDeleted, TestSearchParams and TestDeletedProjects ([303a2f4](303a2f4870b473f3ab081d4bebc76087b53c15c9))
- Assert proxying, fix issue with compensations not running, increase timeouts ([5dff64e](5dff64e3f5d919e6ced108adeca91fcf83724558))
- Cleanup deployments, refactor compensations ([19cc139](19cc13906a3eb467261caa0723aae18d7d08bc95))
- Remove timeouts ([ef64a85](ef64a85ab687f002e5c457e13e3866c230c16bc3))
- Basic fitness tests ([e8e8415](e8e84156965f071f065251329d3e4073a4bbdb6b))
- Fitness test assert proxy ([50aef35](50aef3544a1cf923d7b6b12ae16a8ace112a6a29))
- Add deploy workflow failure test when trying to compose up, fix issues with deploy workflow and add TODOs, refactor to allow proxying per request services ([6b9f623](6b9f6230d3a3e0e44b8ab7c8bb03dc30b2f2c79a))
- Buang deployment errors, fix issues with dbos workflow, fix buang deployment activity return ([1a877bf](1a877bf5c4a3abed84cc3c204e21e98e9a654760))
- Test git clone err ([c5cb017](c5cb017f798060ec33d9c500a680cec34fb92cd1))

<!-- generated by git-cliff -->
