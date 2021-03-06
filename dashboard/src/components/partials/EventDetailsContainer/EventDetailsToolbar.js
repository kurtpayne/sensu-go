import React from "react";
import PropTypes from "prop-types";
import gql from "graphql-tag";

import ClearSilenceAction from "/components/partials/ClearSilenceAction";
import DeleteMenuItem from "/components/partials/ToolbarMenuItems/Delete";
import QueueMenuItem from "/components/partials/ToolbarMenuItems/QueueExecution";
import ResolveMenuItem from "/components/partials/ToolbarMenuItems/Resolve";
import SilenceMenuItem from "/components/partials/ToolbarMenuItems/Silence";
import UnsilenceMenuItem from "/components/partials/ToolbarMenuItems/Unsilence";
import Toolbar from "/components/partials/Toolbar";
import ToolbarMenu from "/components/partials/ToolbarMenu";

import DeleteAction from "./EventDetailsDeleteAction";
import ResolveAction from "./EventDetailsResolveAction";
import ReRunAction from "./EventDetailsReRunAction";
import SilenceAction from "./EventDetailsSilenceAction";

class EventDetailsToolbar extends React.Component {
  static propTypes = {
    event: PropTypes.object.isRequired,
    refetch: PropTypes.func.isRequired,
  };

  static defaultProps = {
    refetch: () => null,
  };

  static fragments = {
    event: gql`
      fragment EventDetailsToolbar_event on Event {
        ...EventDetailsDeleteAction_event
        ...EventDetailsResolveAction_event
        ...EventDetailsReRunAction_event
        ...EventDetailsSilenceAction_event
        ...ClearSilenceAction_record
        isSilenced
      }

      ${DeleteAction.fragments.event}
      ${ResolveAction.fragments.event}
      ${ReRunAction.fragments.event}
      ${SilenceAction.fragments.event}
      ${ClearSilenceAction.fragments.record}
    `,
  };

  render() {
    const { event, refetch } = this.props;

    return (
      <Toolbar
        right={
          <ToolbarMenu fillWidth>
            <ToolbarMenu.Item id="resolve" visible="always">
              <ResolveAction event={event}>
                {({ resolve, canResolve }) => (
                  <ResolveMenuItem disabled={!canResolve} onClick={resolve} />
                )}
              </ResolveAction>
            </ToolbarMenu.Item>
            <ToolbarMenu.Item id="re-run" visible="if-room">
              {event.check.name !== "keepalive" && (
                <ReRunAction event={event}>
                  {run => (
                    <QueueMenuItem
                      title="Re-run Check"
                      titleCondensed="Re-run"
                      onClick={run}
                    />
                  )}
                </ReRunAction>
              )}
            </ToolbarMenu.Item>
            <ToolbarMenu.Item
              id="silence"
              visible={event.isSilenced ? "never" : "if-room"}
            >
              <SilenceAction event={event} onDone={refetch}>
                {menu => <SilenceMenuItem onClick={menu.open} />}
              </SilenceAction>
            </ToolbarMenu.Item>
            <ToolbarMenu.Item
              id="unsilence"
              visible={event.isSilenced ? "if-room" : "never"}
            >
              <ClearSilenceAction record={event} onDone={refetch}>
                {menu => (
                  <UnsilenceMenuItem
                    onClick={menu.open}
                    disabled={!menu.canOpen}
                  />
                )}
              </ClearSilenceAction>
            </ToolbarMenu.Item>
            <ToolbarMenu.Item id="delete" visible="never">
              <DeleteAction event={event}>
                {handler => <DeleteMenuItem onClick={handler} />}
              </DeleteAction>
            </ToolbarMenu.Item>
          </ToolbarMenu>
        }
      />
    );
  }
}

export default EventDetailsToolbar;
