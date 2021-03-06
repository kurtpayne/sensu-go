// +build integration

package schedulerd

import (
	"context"
	"testing"

	"github.com/sensu/sensu-go/backend/messaging"
	"github.com/sensu/sensu-go/backend/store"
	"github.com/sensu/sensu-go/testing/mockring"
	"github.com/sensu/sensu-go/testing/mockstore"
	"github.com/sensu/sensu-go/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCheckWatcherSmoke(t *testing.T) {
	st := &mockstore.MockStore{}
	ringGetter := &mockring.Getter{}

	bus, err := messaging.NewWizardBus(messaging.WizardBusConfig{
		RingGetter: ringGetter,
	})
	require.NoError(t, err)
	require.NoError(t, bus.Start())
	defer bus.Stop()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	checkA := types.FixtureCheckConfig("a")
	checkB := types.FixtureCheckConfig("b")
	st.On("GetCheckConfigs", mock.Anything).Return([]*types.CheckConfig{checkA, checkB}, nil)
	st.On("GetCheckConfigByName", mock.Anything, "a").Return(checkA, nil)
	st.On("GetCheckConfigByName", mock.Anything, "b").Return(checkB, nil)
	st.On("GetAssets", mock.Anything).Return([]*types.Asset{}, nil)
	st.On("GetHookConfigs", mock.Anything).Return([]*types.HookConfig{}, nil)

	watcherChan := make(chan store.WatchEventCheckConfig)
	st.On("GetCheckConfigWatcher", mock.Anything).Return((<-chan store.WatchEventCheckConfig)(watcherChan), nil)

	watcher := NewCheckWatcher(bus, st, ctx)
	require.NoError(t, watcher.Start())

	checkA.Interval = 5
	watcherChan <- store.WatchEventCheckConfig{
		CheckConfig: checkA,
		Action:      store.WatchUpdate,
	}

	checkB.Cron = "* * * * *"
	watcherChan <- store.WatchEventCheckConfig{
		CheckConfig: checkB,
		Action:      store.WatchUpdate,
	}

	watcherChan <- store.WatchEventCheckConfig{
		CheckConfig: checkA,
		Action:      store.WatchDelete,
	}

	watcherChan <- store.WatchEventCheckConfig{
		CheckConfig: checkB,
		Action:      store.WatchCreate,
	}
}
