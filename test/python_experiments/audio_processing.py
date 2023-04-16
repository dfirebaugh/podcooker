import argparse
from pydub import AudioSegment, effects
from concurrent.futures import ThreadPoolExecutor

def convert_to_wav(file_path):
    if not file_path.endswith('.wav'):
        sound = AudioSegment.from_file(file_path)
        file_path = file_path.replace('.mp3', '.wav')
        sound.export(file_path, format='wav')
    return file_path

def normalize_audio(sound, target_dBFS=-22):
    """
    Normalize the audio to a specific dBFS level.
    """
    change_in_dBFS = target_dBFS - sound.dBFS
    return sound.apply_gain(change_in_dBFS)

def compress_audio(sound, threshold=-20.0, ratio=6.0, attack=5, release=50):
    """
    Apply dynamic range compression to the audio.
    """
    return effects.compress_dynamic_range(sound, threshold=threshold, ratio=ratio, attack=attack, release=release)

def apply_limiter(sound, threshold=-1.0):
    """
    Apply a limiter to the audio.
    """
    return effects.limit(sound, threshold)

def apply_gate(sound, threshold=-40, attack=100, release=100):
    """
    Apply a noise gate to the audio.
    """
    return effects.noise_gate(sound, threshold, attack, release)

def process_track(track):
    processed = AudioSegment.from_file(track)
    # processed = compress_audio(processed)
    # processed = apply_limiter(processed)
    # processed = apply_gate(processed)
    # processed = compress_audio(processed)
    processed = normalize_audio(processed)
    return processed

def load_input_tracks(tracks):
    inputs = []
    with ThreadPoolExecutor() as executor:
        inputs = list(executor.map(process_track, tracks))
    return inputs

def mix_tracks(inputs):
    mixed = inputs[0]
    for i in range(1, len(inputs)):
        mixed = mixed.overlay(inputs[i], position=0)
    return mixed

def load_intro_outro_tracks(intro_track, outro_track):
    intro = None
    if intro_track:
        if intro_track.endswith('.mp3'):
            intro_track = convert_to_wav(intro_track)
        intro = AudioSegment.from_file(intro_track)[:30000].apply_gain(-15).fade_out(5000)

    outro = None
    if outro_track:
        if outro_track.endswith('.mp3'):
            outro_track = convert_to_wav(outro_track)
        outro = AudioSegment.from_file(outro_track)[-30000:].apply_gain(-15).fade_in(5000)
        outro = outro.append(AudioSegment.silent(duration=30000))

    return intro, outro

def add_intro(mixed, intro):
    if intro:
        mixed = intro.append(mixed)
    return mixed

def add_outro(mixed, outro):
    if outro:
        mixed = mixed.append(outro.fade_out(len(outro) // 2))
    return mixed

def strip_silence(sound):
    """
    Remove silence from the beginning and end of the audio.
    """
    return sound.strip_silence(silence_len=300, silence_thresh=-40)

def export_audio(sound, output_file='final_output.mp3'):
    """
    Export mixed and stripped audio to an mp3 file.
    """
    sound.export(output_file, format='mp3')
    print(f'Done! Saved to {output_file}')

if __name__ == '__main__':
    # Parsing command line arguments
    parser = argparse.ArgumentParser(description='Process audio files.')
    parser.add_argument('--input_files', nargs='+', help='List of input audio files')
    parser.add_argument('--intro_track', type=str, help='Path to intro track file')
    parser.add_argument('--outro_track', type=str, help='Path to outro track file')
    parser.add_argument('--output_file', type=str, default='final_output.mp3', help='Output audio file')

    args = parser.parse_args()

    # Convert files to wav format
    print('Converting input files to WAV format...')
    input_files = [convert_to_wav(f) for f in args.input_files]

    # Load input tracks and apply effects
    print('Loading input tracks and applying effects...')
    inputs = load_input_tracks(input_files)

    # Mix input tracks together
    print('Mixing input tracks together...')
    mixed_audio = mix_tracks(inputs)

    # Load intro and outro tracks
    print('Loading intro and outro tracks...')
    intro, outro = load_intro_outro_tracks(args.intro_track, args.outro_track)

    # Add intro and outro tracks if available
    print('Adding intro and outro tracks (if available)...')
    mixed_audio = add_intro(mixed_audio, intro)
    mixed_audio = add_outro(mixed_audio, outro)

    # Remove silence from beginning and end of audio
    print('Removing silence from beginning and end of audio...')
    stripped_audio = strip_silence(mixed_audio)

    # Export mixed and stripped audio to an mp3 file
    print(f'Exporting audio to {args.output_file}...')
    export_audio(stripped_audio, args.output_file)

    print('Done!')
