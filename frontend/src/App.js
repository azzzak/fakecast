import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { makeStyles } from '@material-ui/core/styles';
import Drawer from '@material-ui/core/Drawer';
import Box from '@material-ui/core/Box';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemText from '@material-ui/core/ListItemText';
import IconButton from '@material-ui/core/IconButton';
import AddCircleOutlineIcon from '@material-ui/icons/AddCircleOutline';

import Channel from './Channel';

var host;
if (!process.env.NODE_ENV || process.env.NODE_ENV === 'development') {
  host = 'http://127.0.0.1:3333';
} else {
  host = host = '.';
}

const drawerWidth = 240;

const useStyles = makeStyles((theme) => ({
  root: {
    display: 'flex',
  },
  appBar: {
    width: `calc(100% - ${drawerWidth}px)`,
    marginLeft: drawerWidth,
  },
  drawer: {
    width: drawerWidth,
    flexShrink: 0,
  },
  drawerPaper: {
    width: drawerWidth,
  },
  listItemText: {
    fontSize: '1.2em',
  },
  title: {
    whiteSpace: 'nowrap',
    overflow: 'hidden',
    textOverflow: 'ellipsis',
  },
  content: {
    flexGrow: 1,
    padding: theme.spacing(2),
  },
  cont: {
    textAlign: 'center',
    marginBottom: theme.spacing(2),
  },
  placeholder: {
    color: '#555',
    fontSize: '1.2em',
    marginTop: '16px'
  }
}));

export default function PermanentDrawerLeft() {
  const classes = useStyles();

  const [data, setData] = useState([]);
  const [current, setCurrent] = useState({});

  useEffect(() => {
    const fetchData = async () => {
      const res = await axios.get(`${host}/api/list`, {});
      if (res.data) {
        setData(res.data);
      }
    };
    fetchData();
  }, []);

  const addChannel = async () => {
    const res = await axios.post(`${host}/api/channel`);
    const newData = [
      ...data,
      res.data
    ];
    setData(newData);
    setCurrent(res.data);
  };

  const modifyChannel = (channel) => {
    if (!channel.title && !channel.alias) {
      setData(data.filter((item) => item.id === channel.id ? false : true));
      setCurrent({});
      return
    }
    setData(data.map((item) => item.id === channel.id ? channel : item));
  };

  return (
    <div className={classes.root}>
      <Drawer
        className={classes.drawer}
        variant="permanent"
        classes={{
          paper: classes.drawerPaper,
        }}
        anchor="left"
      >
        <List component="nav"
          aria-label="main mailbox folders"
          className={classes.list}
        >
          {data.length > 0 ? data.map((row, index) => (
            row ?
              <ListItem button
                key={row.id}
                className={classes.title}
                onClick={() => setCurrent(data[index])}>
                <ListItemText
                  classes={{ primary: classes.listItemText }}
                  primary={row.title}
                />
              </ListItem> : ''
          )) : <Box p={0}
            display="flex"
            align="center"
            justifyContent="center"
            className={classes.placeholder}
          >
              <Box p={0}>
                Click ⊕ to Add<br />a Channel
             </Box>
            </Box>}
        </List>
        {/* <Divider /> */}

        <div className={classes.cont}>
          <IconButton
            color="primary"
            className={classes.button}
            aria-label="add"
            onClick={() => addChannel()}
          >
            <AddCircleOutlineIcon fontSize="large" />
          </IconButton>
        </div>

      </Drawer>
      <main className={classes.content}>
        <Channel channel={current.id}
          updater={modifyChannel} />
      </main>
    </div >
  );
}