import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { makeStyles } from '@material-ui/core/styles';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import Switch from '@material-ui/core/Switch';
import Button from '@material-ui/core/Button';
import SaveIcon from '@material-ui/icons/Save';
import CancelIcon from '@material-ui/icons/Cancel';
import DeleteIcon from '@material-ui/icons/Delete';
import Drawer from '@material-ui/core/Drawer';
import Grid from '@material-ui/core/Grid';
import Container from '@material-ui/core/Container';

import { field, upd, trimmer } from './Utils';

var host;
if (!process.env.NODE_ENV || process.env.NODE_ENV === 'development') {
  host = 'http://127.0.0.1:3333';
} else {
  host = host = '.';
}

const useStyles = makeStyles(theme => ({
  form: {
    '& .MuiTextField-root': {
      width: '100%'
    },
  },
  button: {
    margin: theme.spacing(1),
  },
  drawer: {
    padding: theme.spacing(2),
  },
}));

export default function Details(props) {
  const [details, setDetails] = useState({});
  const [duration, setDuration] = useState({ h: 0, m: 0, s: 0 });
  const [saved, setSaved] = useState(false);
  const [error, setError] = useState({});

  const classes = useStyles();

  useEffect(() => {
    (async () => {
      if (props.podcast) {
        const result = await axios.get(`${host}/api/channel/${props.channel}/podcast/${props.podcast}`);
        setDetails(result.data);
        setSaved(!result.data.published);
        initDuration(result.data.duration || 0);
      }
    })();
  }, [props.channel, props.podcast]);

  const initDuration = (d) => {
    d = Number(d);
    var h = Math.floor(d / 3600);
    var m = Math.floor(d % 3600 / 60);
    var s = Math.floor(d % 3600 % 60);
    setDuration({ h: h, m: m, s: s })
  };

  const updateRoutine = (method, cb) => {
    props.setDrawer(!props.drawer);
    setSaved(false);
    setError({});
    props.setPodcast();
    props.setPodcasts(method.call([...props.podcasts], cb));
  };

  const updatePodcast = async () => {
    const trimDetails = trimmer(details);
    const { h, m, s } = duration;
    trimDetails.duration = parseInt(h * 3600 + m * 60 + s, 10);
    trimDetails.published = 1;
    await axios.put(`${host}/api/channel/${props.channel}/podcast/${details.id}`, trimDetails);
    setDetails(trimDetails);
    updateRoutine(
      Array.prototype.map,
      (item) => item.id === props.podcast ? details : item
    );
  };

  const deletePodcast = async () => {
    await axios.delete(`${host}/api/channel/${props.channel}/podcast/${details.id}`);
    setDetails({});
    updateRoutine(
      Array.prototype.filter,
      (item) => item.id === props.podcast ? false : true
    );
  };

  const msg = {
    'title': 'Title must not be empty',
  };

  const up = upd(msg, details, error);

  const update = (event, p) => {
    const { updated, err, changed } = up(event, p);
    setDetails(updated);
    setSaved(changed);
    setError(err);
  };

  const numberUpdate = (event, p) => {
    const v = event.target.value;
    var n = parseInt(v, 10);
    n = n >= 0 ? n : 0;
    setDetails({ ...details, [p]: n });
    setSaved(true);
  };

  const handleDuration = (event, p) => {
    const v = event.target.value;
    var n = parseInt(v, 10);
    n = n >= 0 && n < 60 ? n : 0;
    setDuration({ ...duration, [p]: n });
    setSaved(true);
  };

  const handleExplicit = (event) => {
    setDetails({ ...details, 'explicit': +event.target.checked });
    setSaved(true);
  };

  return (
    <Drawer
      classes={{
        paper: classes.drawer,
      }}
      open={props.drawer}
      anchor="bottom"
    >
      <Container maxWidth="md" style={{ marginTop: '8px' }}>
        <form className={classes.form} style={{
          minWidth: '920px',
        }}>

          <Grid container
            spacing={2}
            direction="column"
          >

            <Grid item xs={12}>
              <Grid container direction="row" spacing={2}>
                <Grid item xs={9}>
                  {field({
                    label: "Title",
                    error: !!error.title,
                    helperText: error.title,
                    value: details.title || '',
                    onChange: (e) => update(e, 'title')
                  })}
                </Grid>

                <Grid item xs>
                  {field({
                    label: "Season",
                    type: "number",
                    value: details.season || 0,
                    onChange: (e) => numberUpdate(e, 'season')
                  })}
                </Grid>

                <Grid item xs>
                  {field({
                    label: "Episode",
                    type: "number",
                    value: details.episode || 0,
                    onChange: (e) => numberUpdate(e, 'episode')
                  })}
                </Grid>
              </Grid>
            </Grid>

            <Grid item xs={12}>
              <Grid container
                direction="row"
                spacing={2}
              >
                <Grid item xs={12}>
                  {field({
                    label: "Description",
                    value: details.description || '',
                    rows: 2,
                    onChange: (e) => update(e, 'description')
                  })}
                </Grid>
              </Grid>
            </Grid>

            <Grid item xs={12}>
              <Grid container
                direction="row"
                spacing={2}
              >
                <Grid item xs={8}><FormControlLabel
                  control={
                    <Switch
                      checked={!!details.explicit}
                      onChange={handleExplicit}
                      name="checkedB"
                      color="primary"
                    />
                  }
                  label="Explicit"
                />
                </Grid>

                <Grid item xs>
                  {field({
                    label: "Hours",
                    type: "number",
                    value: duration.h || 0,
                    onChange: (e) => handleDuration(e, 'h')
                  })}
                </Grid>

                <Grid item xs>
                  {field({
                    label: "Minutes",
                    type: "number",
                    value: duration.m || 0,
                    onChange: (e) => handleDuration(e, 'm')
                  })}
                </Grid>

                <Grid item xs>
                  {field({
                    label: "Seconds",
                    type: "number",
                    value: duration.s || 0,
                    onChange: (e) => handleDuration(e, 's')
                  })}
                </Grid>
              </Grid>
            </Grid>

            <Grid item xs={12}>
              <Grid container
                direction="row"
                justify="space-between"
                spacing={0}
              >
                <Grid item>
                  {details.published ? <Button
                    disabled={saved ? false : true}
                    onClick={() => updatePodcast()}
                    variant="outlined"
                    color="primary"
                    size="small"
                    className={classes.button}
                    startIcon={<SaveIcon />}
                  >Save</Button> : <Button
                    disabled={saved ? false : true}
                    onClick={() => updatePodcast()}
                    variant="outlined"
                    color="primary"
                    size="small"
                    className={classes.button}
                    startIcon={<SaveIcon />}
                  >Publish</Button>}
                  <Button
                    onClick={() => {
                      props.setPodcast();
                      props.setDrawer(!props.drawer);
                      setError({})
                    }}
                    variant="outlined"
                    size="small"
                    className={classes.button}
                    startIcon={<CancelIcon />}
                  >Cancel</Button>
                </Grid><Grid item>
                  <Button
                    onClick={() => deletePodcast()}
                    color="secondary"
                    size="small"
                    className={classes.button}
                    startIcon={<DeleteIcon />}
                  >Delete</Button>
                </Grid>
              </Grid>
            </Grid>
          </Grid>
        </form>
      </Container>
    </Drawer >
  );
}