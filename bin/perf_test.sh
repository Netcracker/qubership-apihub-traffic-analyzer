#!/bin/bash
if ! [ -x $1 ]
then
  echo "$1 not an executable"
  exit 1
fi

PERL_CMD=`which perl 2>/dev/null`
if ! [ -x ${PERL_CMD} ]
then
  echo "Perl required to run"
  exit 1
fi
PERF_CMD=`which perf 2>/dev/null`
if ! [ -x ${PERL_CMD} ]
then
  echo "Perf required to run"
  exit 1
fi
FILE_PREFIX="PERF"
FG_PATH=/mnt/c/Projects/unsorted/GitHub/FlameGraph
if ! [ -d ${FG_PATH} ]
then
  echo "Flame graph scripts required to run"
  exit 1
fi
if ! [ -x "${FG_PATH}/stackcollapse-perf.pl" ]
then
  echo "${FG_PATH}/stackcollapse-perf.pl not an executable"
  exit 1
fi
if ! [ -x "${FG_PATH}/flamegraph.pl" ]
then
  echo "${FG_PATH}/flamegraph.pl not an executable"
  exit 1
fi
ENV_FILES_PATH=/mnt/c/Projects/GIT/qubership-apihub-traffic-analyzer
. ${ENV_FILES_PATH}/Service.env
nohup $1 &
PID=`echo $!`
FILE_NAME="${FILE_PREFIX}_${PID}_AHTA"
TPID=`ps -aef | grep ${PID} | grep -v grep`
if [ -z "${TPID}" ]
then
  echo "process $PID has been dead"
  exit 2
fi
sudo $PERF_CMD record -F 99 -p ${PID} -g -o "${FILE_NAME}.perf"
if [ -r "${FILE_PREFIX}_${PID}_AHTA.perf" ]
then
  sudo chmod 0644 "${FILE_NAME}.perf"
  $PERF_CMD script -i "${FILE_NAME}.perf" > "${FILE_NAME}.txt"
  if [ -r "${FILE_NAME}.txt" ]
  then
    cat "${FILE_NAME}.txt" | $FG_PATH/stackcollapse-perf.pl > "${FILE_NAME}.collapsed" && \
    $FG_PATH/flamegraph.pl "${FILE_NAME}.collapsed" > "${FILE_NAME}.svg"
  else
    echo "File ${FILE_NAME}.txt is not readable"
  fi
else
  echo "File ${FILE_NAME}.perf is not readable"
fi
