import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { makeStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import Container from '@material-ui/core/Container';
import Button from '@material-ui/core/Button';
import DeleteIcon from '@material-ui/icons/Delete';
import SaveIcon from '@material-ui/icons/Save';
import CloudUploadIcon from '@material-ui/icons/CloudUpload';

import Details from './Details';
import Cover from './Cover';
import EnhancedTable from './Table';
import { field, upd, trimmer } from './Utils';

var host;
if (!process.env.NODE_ENV || process.env.NODE_ENV === 'development') {
  host = 'http://127.0.0.1:3333';
} else {
  host = host = '.';
}

const useStyles = makeStyles(theme => ({
  root: {
    backgroundColor: theme.palette.background.paper,
  },
  title: {
    fontSize: '3em',
  },
  link: {
    display: 'inline-block',
    marginTop: '0.5em',
    marginBottom: '1.5em',
    color: '#909497',
    '&:hover': {
      color: '#676767',
    },
  },
  form: {
    '& .MuiTextField-root': {
      width: '100%'
    },
  },
  leftInput: {
    width: [theme.spacing(32), '!important'],
    margin: theme.spacing(1),
  },
  button: {
    margin: theme.spacing(1),
  },
  upload: {
    display: 'none',
  },
}));

export default function Content(props) {
  const [info, setInfo] = useState({});
  const [alias, setAlias] = useState('');
  const [changed, setChanged] = useState(false);
  const [drawer, setDrawer] = useState(false);
  const [podcast, setPodcast] = useState();
  const [podcasts, setPodcasts] = useState([]);
  const [error, setError] = useState({});

  const classes = useStyles();

  useEffect(() => {
    (async () => {
      if (props.channel) {
        const res = await axios.get(`${host}/api/channel/${props.channel}`, {});
        setInfo(res.data.info);
        setAlias(res.data.info.alias);
        if (res.data.podcasts)
          setPodcasts(res.data.podcasts);
        else
          setPodcasts([]);
        setError({});
        setChanged(false);
      }
    })();
  }, [props.channel]);

  const updateChannel = async () => {
    const trimInfo = trimmer(info);
    setInfo(trimInfo);
    const [cvr] = trimInfo.cover.split('/').slice(-1);
    const infoUpd = {
      channel: { ...trimInfo, 'cover': cvr },
      old_alias: alias
    }
    const res = await axios.put(`${host}/api/channel/${info.id}`, infoUpd);

    if (res.data.error) {
      setError({ ...error, 'alias': errInUse });
      return
    }

    props.updater({
      id: trimInfo.id,
      alias: trimInfo.alias,
      title: trimInfo.title
    });

    if (alias !== trimInfo.alias) {
      setAlias(trimInfo.alias);
      setInfo({ ...info, 'cover': res.data.cover });
    }

    setChanged(false);
  };

  const uploadPodcast = async (e) => {
    var formData = new FormData();
    const podcastFile = e.target.files[0];
    formData.append('file', podcastFile);
    formData.append('length', podcastFile.size);
    const res = await axios.post(`${host}/api/channel/${info.id}/upload`, formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    });
    setPodcasts([
      res.data,
      ...podcasts
    ]);
    setPodcast(res.data.id);
    setDrawer(true);
  };

  const uploadCover = async (e) => {
    var formData = new FormData();
    formData.append('file', e.target.files[0]);
    const res = await axios.post(`${host}/api/channel/${info.id}/cover/upload`, formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    });
    setInfo({ ...info, 'cover': res.data.cover });
  };

  const deleteChannel = async () => {
    props.updater({ id: info.id });
    await axios.delete(`${host}/api/channel/${info.id}`);
    setInfo({});
    setPodcasts([]);
    setChanged(false);
  };

  const deleteCover = async () => {
    const fields = info.cover.split('/');
    const [, , cover] = fields.slice(-3);
    await axios.delete(`${host}/api/channel/${info.id}/cover/${cover}`);
    setInfo({ ...trimmer(info), 'cover': '' });
  };

  const errEmpty = {
    'title': 'Title must not be empty',
    'alias': 'Alias must not be empty',
  };

  const errInvalid = 'Alias has incorrect symbols';
  const errInUse = "Can't use this name, it may be already in use";

  const up = upd(errEmpty, info, error);

  const update = (event, p) => {
    const { updated, err, changed } = up(event, p);
    setInfo(updated);
    setChanged(changed);
    setError(err);
  };

  const handleAlias = (e) => {
    update(e, 'alias');
    const re = /^([a-z0-9_-])*$/;
    if (!re.test(e.target.value)) {
      setError({ ...error, 'alias': errInvalid });
      setChanged(false);
    }
  };

  return (
    <div className={classes.root}>
      {info.id ?
        <Container maxWidth="lg" style={{ minWidth: '720px', margin: '20px 0 10px 0' }}>
          <form className={classes.form}>
            <Grid container
              spacing={2}
              direction="row"
            >
              <Grid style={{ width: '294px' }}>
                <Grid item>
                  <Grid>
                    <Cover
                      cover={info.cover}
                      upload={uploadCover}
                      delete={deleteCover}
                    />
                  </Grid>
                  <Grid style={{ marginTop: '30px' }}>
                    {field({
                      label: "Author",
                      className: classes.leftInput,
                      value: info.author || '',
                      onChange: (e) => update(e, 'author')
                    })}
                  </Grid>
                </Grid>
              </Grid>
              <Grid item xs>
                <Grid container
                  direction="column"
                >
                  <Grid item xs>
                    <div className={classes.title}>{info.title}</div>
                    <a href={`${info.host}/feed/${info.alias}`} className={classes.link} >{`${info.host}/feed/${info.alias}`}</a>
                  </Grid>

                  <Grid container
                    direction="column"
                    spacing={2}
                  >
                    <Grid item xs>

                      <Grid container
                        direction="row"
                        spacing={2}
                      >
                        <Grid item xs={8}>
                          {field({
                            label: "Title",
                            error: !!error.title,
                            helperText: error.title,
                            value: info.title || '',
                            onChange: (e) => update(e, 'title')
                          })}
                        </Grid>
                        <Grid item xs>
                          {field({
                            label: "Alias",
                            error: !!error.alias,
                            helperText: error.alias,
                            value: info.alias || '',
                            onChange: (e) => handleAlias(e)
                          })}
                        </Grid>
                      </Grid>
                    </Grid>

                    <Grid item xs>
                      {field({
                        label: "Description",
                        value: info.description || '',
                        rows: 2,
                        onChange: (e) => update(e, 'description')
                      })}
                    </Grid>
                  </Grid>

                  <Grid>
                    <Grid container
                      direction="row"
                      justify="space-between"
                      style={{ marginTop: '10px' }}
                    >
                      <Grid item>
                        <input
                          accept="audio/x-m4a, audio/mpeg"
                          className={classes.upload}
                          id="file"
                          name="file"
                          type="file"
                          onChange={(e) => {
                            uploadPodcast(e);
                            e.target.value = null;
                          }}
                        />
                        <label htmlFor="file">
                          <Button
                            variant="outlined"
                            color="primary"
                            size="small"
                            className={classes.button}
                            component="span"
                            startIcon={<CloudUploadIcon />}
                          >Upload</Button></label>
                        <Button
                          disabled={changed ? false : true}
                          onClick={() => updateChannel()}
                          variant="outlined"
                          color="primary"
                          size="small"
                          className={classes.button}
                          startIcon={<SaveIcon />}
                        >Save</Button>
                      </Grid>
                      <Grid item>
                        <Button
                          onClick={() => deleteChannel()}
                          color="secondary"
                          size="small"
                          className={classes.button}
                          startIcon={<DeleteIcon />}
                        >Delete</Button>
                      </Grid>
                    </Grid>
                  </Grid>
                  <Grid style={{ margin: '16px 0' }}>
                    {podcasts.length > 0 ? <EnhancedTable podcasts={podcasts}
                      setDrawer={setDrawer}
                      setPodcast={setPodcast} /> : ''}
                  </Grid>
                </Grid>
              </Grid>
            </Grid>
          </form ></Container> : ''
      }
      <Details
        channel={info.id}
        podcast={podcast}
        setPodcast={setPodcast}
        drawer={drawer}
        setDrawer={setDrawer}
        podcasts={podcasts}
        setPodcasts={setPodcasts}
      />
    </div >
  );
}