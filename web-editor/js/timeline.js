// 时间轴控制系统 - timeline.js

// 关键帧类
class Keyframe {
    constructor(time, property, value, ease = 'linear') {
        this.id = Utils.generateId();
        this.time = time; // 时间（秒）
        this.property = property; // 属性名
        this.value = value; // 属性值
        this.ease = ease; // 缓动类型
    }

    clone() {
        return new Keyframe(this.time, this.property, this.value, this.ease);
    }
}

// 动画轨道类
class Track {
    constructor(objectId, property) {
        this.id = Utils.generateId();
        this.objectId = objectId;
        this.property = property;
        this.keyframes = [];
        this.visible = true;
        this.locked = false;
    }

    addKeyframe(time, value, ease = 'linear') {
        // 如果在相同时间已有关键帧，则替换
        const existingIndex = this.keyframes.findIndex(kf => kf.time === time);
        const keyframe = new Keyframe(time, this.property, value, ease);

        if (existingIndex >= 0) {
            this.keyframes[existingIndex] = keyframe;
        } else {
            this.keyframes.push(keyframe);
            this.sortKeyframes();
        }

        return keyframe;
    }

    removeKeyframe(keyframeId) {
        this.keyframes = this.keyframes.filter(kf => kf.id !== keyframeId);
    }

    getKeyframe(keyframeId) {
        return this.keyframes.find(kf => kf.id === keyframeId);
    }

    getKeyframeAt(time) {
        return this.keyframes.find(kf => kf.time === time);
    }

    sortKeyframes() {
        this.keyframes.sort((a, b) => a.time - b.time);
    }

    // 获取指定时间的插值
    getValueAt(time) {
        if (this.keyframes.length === 0) return null;
        if (this.keyframes.length === 1) return this.keyframes[0].value;

        // 查找关键帧区间
        let startKf = null, endKf = null;

        for (let i = 0; i < this.keyframes.length; i++) {
            const kf = this.keyframes[i];
            if (kf.time <= time) {
                startKf = kf;
            }
            if (kf.time >= time && !endKf) {
                endKf = kf;
                break;
            }
        }

        if (!startKf) return this.keyframes[0].value;
        if (!endKf) return this.keyframes[this.keyframes.length - 1].value;
        if (startKf === endKf) return startKf.value;

        // 插值计算
        const t = (time - startKf.time) / (endKf.time - startKf.time);
        const easedT = this.applyEase(t, endKf.ease);

        return this.interpolateValue(startKf.value, endKf.value, easedT);
    }

    interpolateValue(from, to, t) {
        if (typeof from === 'number' && typeof to === 'number') {
            return from + (to - from) * t;
        }

        if (typeof from === 'string' && typeof to === 'string') {
            // 颜色插值
            if (from.startsWith('#') && to.startsWith('#')) {
                return this.interpolateColor(from, to, t);
            }
            // 其他字符串直接切换
            return t < 0.5 ? from : to;
        }

        // 对象插值（如位置、大小等）
        if (typeof from === 'object' && typeof to === 'object') {
            const result = {};
            for (let key in from) {
                if (key in to) {
                    result[key] = this.interpolateValue(from[key], to[key], t);
                } else {
                    result[key] = from[key];
                }
            }
            return result;
        }

        return t < 0.5 ? from : to;
    }

    interpolateColor(from, to, t) {
        const fromRgb = Utils.hexToRgb(from);
        const toRgb = Utils.hexToRgb(to);

        const r = Math.round(fromRgb.r + (toRgb.r - fromRgb.r) * t);
        const g = Math.round(fromRgb.g + (toRgb.g - fromRgb.g) * t);
        const b = Math.round(fromRgb.b + (toRgb.b - fromRgb.b) * t);

        return Utils.rgbToHex(r, g, b);
    }

    applyEase(t, ease) {
        switch (ease) {
            case 'linear':
                return t;
            case 'ease-in':
                return t * t;
            case 'ease-out':
                return 1 - (1 - t) * (1 - t);
            case 'ease-in-out':
                return t < 0.5 ? 2 * t * t : 1 - 2 * (1 - t) * (1 - t);
            case 'bounce':
                if (t < 0.5) {
                    return 2 * t * t;
                } else {
                    return 1 - 2 * (1 - t) * (1 - t);
                }
            default:
                return t;
        }
    }

    toJSON() {
        return {
            id: this.id,
            objectId: this.objectId,
            property: this.property,
            keyframes: this.keyframes,
            visible: this.visible,
            locked: this.locked
        };
    }

    static fromJSON(data) {
        const track = new Track(data.objectId, data.property);
        track.id = data.id;
        track.visible = data.visible;
        track.locked = data.locked;
        track.keyframes = data.keyframes.map(kfData => {
            const kf = new Keyframe(kfData.time, kfData.property, kfData.value, kfData.ease);
            kf.id = kfData.id;
            return kf;
        });
        return track;
    }
}

// 时间轴管理器
class Timeline extends EventEmitter {
    constructor() {
        super();
        this.currentTime = 0;
        this.duration = 5; // 总时长（秒）
        this.frameRate = 30; // 帧率
        this.tracks = [];
        this.selectedKeyframes = [];
        this.isPlaying = false;
        this.playStartTime = 0;
        this.animationFrame = null;

        // UI元素
        this.container = null;
        this.timelineContent = null;
        this.scrubber = null;
        this.timeDisplay = null;

        // 配置
        this.pixelsPerSecond = 100;
        this.trackHeight = 40;
        this.snapToFrames = true;

        this.init();
    }

    init() {
        this.container = document.getElementById('timelineContainer');
        this.timelineContent = document.getElementById('timelineContent');
        this.scrubber = document.getElementById('timeScrubber');
        this.timeDisplay = document.getElementById('currentTime');

        if (this.container) {
            this.setupEventListeners();
            this.updateTimelineDisplay();
        }
    }

    setupEventListeners() {
        // 播放控制
        const playBtn = document.getElementById('playBtn');
        const pauseBtn = document.getElementById('pauseBtn');
        const stopBtn = document.getElementById('stopBtn');

        if (playBtn) playBtn.addEventListener('click', () => this.play());
        if (pauseBtn) pauseBtn.addEventListener('click', () => this.pause());
        if (stopBtn) stopBtn.addEventListener('click', () => this.stop());

        // 时间轴点击
        if (this.timelineContent) {
            this.timelineContent.addEventListener('click', this.onTimelineClick.bind(this));
            this.timelineContent.addEventListener('mousemove', this.onTimelineMouseMove.bind(this));
        }

        // 拖动控制
        if (this.scrubber) {
            this.scrubber.addEventListener('mousedown', this.startScrubbing.bind(this));
        }

        // 键盘控制
        document.addEventListener('keydown', this.onKeyDown.bind(this));
    }

    // 播放控制
    play() {
        if (this.isPlaying) return;

        this.isPlaying = true;
        this.playStartTime = performance.now() - this.currentTime * 1000;
        this.updatePlayButtons();
        this.animationLoop();
        this.emit('playStateChanged', true);
    }

    pause() {
        this.isPlaying = false;
        if (this.animationFrame) {
            cancelAnimationFrame(this.animationFrame);
            this.animationFrame = null;
        }
        this.updatePlayButtons();
        this.emit('playStateChanged', false);
    }

    stop() {
        this.pause();
        this.setCurrentTime(0);
        this.emit('stopped');
    }

    togglePlayback() {
        if (this.isPlaying) {
            this.pause();
        } else {
            this.play();
        }
    }

    animationLoop() {
        if (!this.isPlaying) return;

        const currentTime = (performance.now() - this.playStartTime) / 1000;

        if (currentTime >= this.duration) {
            this.stop();
            return;
        }

        this.setCurrentTime(currentTime);
        this.animationFrame = requestAnimationFrame(() => this.animationLoop());
    }

    updatePlayButtons() {
        const playBtn = document.getElementById('playBtn');
        const pauseBtn = document.getElementById('pauseBtn');

        if (playBtn) playBtn.style.display = this.isPlaying ? 'none' : 'block';
        if (pauseBtn) pauseBtn.style.display = this.isPlaying ? 'block' : 'none';
    }

    // 时间控制
    setCurrentTime(time) {
        this.currentTime = Math.max(0, Math.min(time, this.duration));
        this.updateTimeDisplay();
        this.updateScrubberPosition();
        this.applyAnimationAtTime(this.currentTime);
        this.emit('timeChanged', this.currentTime);
    }

    setDuration(duration) {
        this.duration = Math.max(0.1, duration);
        this.updateTimelineDisplay();
        this.emit('durationChanged', this.duration);
    }

    // 轨道管理
    createTrackForObject(objectId, property) {
        const existingTrack = this.tracks.find(t =>
            t.objectId === objectId && t.property === property);

        if (existingTrack) {
            return existingTrack;
        }

        const track = new Track(objectId, property);
        this.tracks.push(track);
        this.updateTrackDisplay();
        this.emit('trackAdded', track);

        return track;
    }

    removeTrack(trackId) {
        this.tracks = this.tracks.filter(t => t.id !== trackId);
        this.updateTrackDisplay();
        this.emit('trackRemoved', trackId);
    }

    getTrack(trackId) {
        return this.tracks.find(t => t.id === trackId);
    }

    getTracksForObject(objectId) {
        return this.tracks.filter(t => t.objectId === objectId);
    }

    // 关键帧操作
    addKeyframe(objectId, property, time, value, ease = 'linear') {
        const track = this.createTrackForObject(objectId, property);
        const keyframe = track.addKeyframe(time, value, ease);
        this.updateTrackDisplay();
        this.emit('keyframeAdded', keyframe);
        return keyframe;
    }

    removeKeyframe(trackId, keyframeId) {
        const track = this.getTrack(trackId);
        if (track) {
            track.removeKeyframe(keyframeId);
            this.updateTrackDisplay();
            this.emit('keyframeRemoved', keyframeId);
        }
    }

    // 自动关键帧
    autoKeyframe(objectId) {
        const object = window.app && window.app.objectManager ?
            window.app.objectManager.getObject(objectId) : null;
        if (!object) return;

        const properties = ['x', 'y', 'width', 'height', 'rotation', 'opacity'];
        const animatableProps = properties.filter(prop =>
            object.hasOwnProperty(prop) && typeof object[prop] !== 'function');

        animatableProps.forEach(prop => {
            this.addKeyframe(objectId, prop, this.currentTime, object[prop]);
        });
    }

    // 应用动画
    applyAnimationAtTime(time) {
        this.tracks.forEach(track => {
            const object = window.app && window.app.objectManager ?
                window.app.objectManager.getObject(track.objectId) : null;
            if (!object || !track.visible) return;

            const value = track.getValueAt(time);
            if (value !== null) {
                this.applyPropertyValue(object, track.property, value);
            }
        });
    }

    applyPropertyValue(object, property, value) {
        switch (property) {
            case 'x':
            case 'y':
                object.setPosition(
                    property === 'x' ? value : object.x,
                    property === 'y' ? value : object.y
                );
                break;
            case 'width':
            case 'height':
                object.setSize(
                    property === 'width' ? value : object.width,
                    property === 'height' ? value : object.height
                );
                break;
            case 'rotation':
                object.setRotation(value);
                break;
            case 'opacity':
                object.setOpacity(value);
                break;
            case 'fillColor':
                object.setFillColor(value);
                break;
            case 'strokeColor':
                object.setStrokeColor(value);
                break;
            default:
                if (object[property] !== undefined) {
                    object[property] = value;
                    object.updateElement();
                }
        }
    }

    // UI更新
    updateTimelineDisplay() {
        if (!this.timelineContent) return;

        const width = this.duration * this.pixelsPerSecond;
        this.timelineContent.style.width = width + 'px';

        // 更新时间刻度
        this.updateTimeRuler();
        this.updateTrackDisplay();
    }

    updateTimeRuler() {
        let ruler = document.getElementById('timeRuler');
        if (!ruler) {
            ruler = Utils.createElement('div', 'time-ruler');
            ruler.id = 'timeRuler';
            this.timelineContent.appendChild(ruler);
        }

        ruler.innerHTML = '';
        const step = 1; // 每秒一个刻度

        for (let i = 0; i <= this.duration; i += step) {
            const mark = Utils.createElement('div', 'time-mark');
            mark.style.left = (i * this.pixelsPerSecond) + 'px';
            mark.textContent = i + 's';
            ruler.appendChild(mark);
        }
    }

    updateTrackDisplay() {
        const trackContainer = document.getElementById('trackContainer');
        if (!trackContainer) return;

        trackContainer.innerHTML = '';

        this.tracks.forEach((track, index) => {
            const trackElement = this.createTrackElement(track, index);
            trackContainer.appendChild(trackElement);
        });
    }

    createTrackElement(track, index) {
        const object = window.app && window.app.objectManager ?
            window.app.objectManager.getObject(track.objectId) : null;
        const trackEl = Utils.createElement('div', 'timeline-track');
        trackEl.style.top = (index * this.trackHeight) + 'px';
        trackEl.style.height = this.trackHeight + 'px';

        // 轨道标题
        const title = Utils.createElement('div', 'track-title');
        title.textContent = `${object ? object.name : 'Unknown'} - ${track.property}`;
        trackEl.appendChild(title);

        // 关键帧
        track.keyframes.forEach(keyframe => {
            const kfEl = this.createKeyframeElement(keyframe, track);
            trackEl.appendChild(kfEl);
        });

        return trackEl;
    }

    createKeyframeElement(keyframe, track) {
        const kfEl = Utils.createElement('div', 'keyframe');
        kfEl.style.left = (keyframe.time * this.pixelsPerSecond) + 'px';
        kfEl.dataset.keyframeId = keyframe.id;
        kfEl.dataset.trackId = track.id;

        // 关键帧事件
        kfEl.addEventListener('click', (e) => {
            e.stopPropagation();
            this.selectKeyframe(keyframe.id, e.ctrlKey);
        });

        kfEl.addEventListener('contextmenu', (e) => {
            e.preventDefault();
            this.showKeyframeContextMenu(e.clientX, e.clientY, keyframe, track);
        });

        return kfEl;
    }

    updateTimeDisplay() {
        if (this.timeDisplay) {
            const minutes = Math.floor(this.currentTime / 60);
            const seconds = (this.currentTime % 60).toFixed(2);
            this.timeDisplay.textContent = `${minutes}:${seconds.padStart(5, '0')}`;
        }
    }

    updateScrubberPosition() {
        if (this.scrubber) {
            this.scrubber.style.left = (this.currentTime * this.pixelsPerSecond) + 'px';
        }
    }

    // 事件处理
    onTimelineClick(e) {
        const rect = this.timelineContent.getBoundingClientRect();
        const x = e.clientX - rect.left;
        const time = x / this.pixelsPerSecond;
        this.setCurrentTime(time);
    }

    onTimelineMouseMove(e) {
        // 显示时间提示
        const rect = this.timelineContent.getBoundingClientRect();
        const x = e.clientX - rect.left;
        const time = x / this.pixelsPerSecond;

        // 更新时间提示
        let tooltip = document.getElementById('timeTooltip');
        if (!tooltip) {
            tooltip = Utils.createElement('div', 'time-tooltip');
            tooltip.id = 'timeTooltip';
            document.body.appendChild(tooltip);
        }

        tooltip.textContent = time.toFixed(2) + 's';
        tooltip.style.left = e.clientX + 'px';
        tooltip.style.top = (e.clientY - 30) + 'px';
        tooltip.style.display = 'block';
    }

    startScrubbing(e) {
        e.preventDefault();

        const startX = e.clientX;
        const startTime = this.currentTime;

        const onMouseMove = (e) => {
            const dx = e.clientX - startX;
            const dt = dx / this.pixelsPerSecond;
            this.setCurrentTime(startTime + dt);
        };

        const onMouseUp = () => {
            document.removeEventListener('mousemove', onMouseMove);
            document.removeEventListener('mouseup', onMouseUp);
        };

        document.addEventListener('mousemove', onMouseMove);
        document.addEventListener('mouseup', onMouseUp);
    }

    onKeyDown(e) {
        if (e.target.tagName === 'INPUT') return;

        switch (e.key) {
            case ' ':
                e.preventDefault();
                this.isPlaying ? this.pause() : this.play();
                break;
            case 'Home':
                this.setCurrentTime(0);
                break;
            case 'End':
                this.setCurrentTime(this.duration);
                break;
            case 'ArrowLeft':
                this.setCurrentTime(this.currentTime - 1 / this.frameRate);
                break;
            case 'ArrowRight':
                this.setCurrentTime(this.currentTime + 1 / this.frameRate);
                break;
        }
    }

    // 关键帧选择
    selectKeyframe(keyframeId, addToSelection = false) {
        if (!addToSelection) {
            this.selectedKeyframes = [];
        }

        if (!this.selectedKeyframes.includes(keyframeId)) {
            this.selectedKeyframes.push(keyframeId);
        }

        this.updateKeyframeSelection();
        this.emit('keyframeSelectionChanged', this.selectedKeyframes);
    }

    clearKeyframeSelection() {
        this.selectedKeyframes = [];
        this.updateKeyframeSelection();
    }

    updateKeyframeSelection() {
        document.querySelectorAll('.keyframe').forEach(el => {
            const isSelected = this.selectedKeyframes.includes(el.dataset.keyframeId);
            el.classList.toggle('selected', isSelected);
        });
    }

    // 序列化
    toJSON() {
        return {
            currentTime: this.currentTime,
            duration: this.duration,
            frameRate: this.frameRate,
            tracks: this.tracks.map(track => track.toJSON())
        };
    }

    fromJSON(data) {
        this.currentTime = data.currentTime || 0;
        this.duration = data.duration || 5;
        this.frameRate = data.frameRate || 30;

        this.tracks = (data.tracks || []).map(trackData => Track.fromJSON(trackData));

        this.updateTimelineDisplay();
        this.setCurrentTime(this.currentTime);
    }

    clear() {
        this.tracks = [];
        this.selectedKeyframes = [];
        this.setCurrentTime(0);
        this.updateTrackDisplay();
    }
}

// 导出
window.Timeline = Timeline;
window.Track = Track;
window.Keyframe = Keyframe;
