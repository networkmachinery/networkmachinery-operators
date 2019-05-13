package app

import (
	"context"
	"github.com/networkmachinery/networkmachinery-operators/pkg/apis/networkmachinery/v1alpha1"
	networkconnectivitywebhook "github.com/networkmachinery/networkmachinery-operators/pkg/controllers/networkconnectivity/webhook"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"

	"k8s.io/cli-runtime/pkg/genericclioptions"

	"github.com/networkmachinery/networkmachinery-operators/pkg/controllers"
	admissionregistrationv1beta1 "k8s.io/api/admissionregistration/v1beta1"
	apitypes "k8s.io/apimachinery/pkg/types"

	"github.com/networkmachinery/networkmachinery-operators/pkg/controllers/networkconnectivity/controller"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission/builder"

	"github.com/networkmachinery/networkmachinery-operators/pkg/apis/networkmachinery/install"
	"github.com/networkmachinery/networkmachinery-operators/pkg/utils"
	"github.com/spf13/cobra"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

var log = logf.Log.WithName("example-controller")

const(
	validatingServerName = "networkconnectivity-layer-validator"
	validationServerWebhookSecret = "networkconnectivity-layer-validator-secret"
	validationServerWebhookService = "networkconnectivity-layer-validator-service"
)
func NewNetworkConnectivityTestCmd(ctx context.Context) *cobra.Command {
	entryLog := log.WithName("networkconnectivity-test-cmd")

	networkConnectivityTestCmdOpts := NetworkConnectivityTestCmdOpts{
		ConfigFlags: genericclioptions.NewConfigFlags(),
		LeaderElectionOptions: controllers.LeaderElectionOptions{
			LeaderElection:          true,
			LeaderElectionNamespace: "default",
			LeaderElectionID:        utils.LeaderElectionNameID(controller.Name),
		},
		ControllerOptions: controllers.ControllerOptions{
			MaxConcurrentReconciles: 5,
		},
	}

	cmd := &cobra.Command{
		Use: "networkconnectivity-test-controller",
		Run: func(cmd *cobra.Command, args []string) {
			mgrOptions := &manager.Options{}
			mgr, err := manager.New(networkConnectivityTestCmdOpts.InitConfig(), *networkConnectivityTestCmdOpts.ApplyLeaderElection(mgrOptions))
			if err != nil {
				utils.LogErrAndExit(err, "Could not instantiate manager")
			}
			if err := install.AddToScheme(mgr.GetScheme()); err != nil {
				utils.LogErrAndExit(err, "Could not update	 manager scheme")
			}


			validatingWebhook, err := builder.NewWebhookBuilder().
				Name("validating.k8s.io").
				Validating().
				Operations(admissionregistrationv1beta1.Create, admissionregistrationv1beta1.Update).
				WithManager(mgr).
				ForType(&v1alpha1.NetworkConnectivityTest{}).
				Handlers(&networkconnectivitywebhook.LayerValidator{}).
				Build()
			if err != nil {
				entryLog.Error(err, "unable to setup validating webhook")
				os.Exit(1)
			}

			entryLog.Info("setting up webhook server")
			admissionServer, err := webhook.NewServer(validatingServerName, mgr, webhook.ServerOptions{
				Port:                          9876,
				CertDir:                       "/tmp/cert",
				DisableWebhookConfigInstaller: &networkConnectivityTestCmdOpts.disableWebhookConfigInstaller,
				BootstrapOptions: &webhook.BootstrapOptions{
					Secret: &apitypes.NamespacedName{
						Namespace: metav1.NamespaceDefault,
						Name:      validationServerWebhookSecret,
					},

					Service: &webhook.Service{
						Namespace:  metav1.NamespaceDefault,
						Name:      validationServerWebhookService,
						// Selectors should select the pods that runs this webhook server.
						Selectors: map[string]string{
							"app.kubernetes.io/name":controller.Name,
						},
					},
				},
			})
			if err != nil {
				entryLog.Error(err, "unable to create a new webhook server")
				os.Exit(1)
			}

			entryLog.Info("registering webhooks to the webhook server")
			err = admissionServer.Register(validatingWebhook)
			if err != nil {
				entryLog.Error(err, "unable to register webhooks in the admission server")
				os.Exit(1)
			}

			if err := controllers.AddToManager(mgr); err != nil {
				utils.LogErrAndExit(err, "Could not add controller to manager")
			}

			if err := mgr.Start(ctx.Done()); err != nil {
				utils.LogErrAndExit(err, "Error running manager")
			}
		},
	}

	networkConnectivityTestCmdOpts.AddAllFlags(cmd.Flags())
	return cmd
}
